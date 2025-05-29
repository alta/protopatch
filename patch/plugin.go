package patch

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"slices"
	"strings"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
)

// StripParams strips a named param from req.
func StripParams(req *pluginpb.CodeGeneratorRequest, p []string) {
	if req.Parameter == nil {
		return
	}
	v := stripParams(*req.Parameter, p)
	req.Parameter = &v

}

func stripParams(s string, p []string) string {
	var b strings.Builder
	for _, param := range strings.Split(s, ",") {
		if !slices.Contains(p, strings.SplitN(param, "=", 2)[0]) {
			if b.Len() > 0 {
				b.WriteString(",")
			}
			b.WriteString(param)
		}
	}
	return b.String()
}

// RunPlugin runs a protoc plugin named "protoc-gen-$plugin"
// and returns the generated CodeGeneratorResponse or an error.
// Supply a non-nil stderr to override stderr on the called plugin.
func RunPlugin(plugin string, req *pluginpb.CodeGeneratorRequest, stderr io.Writer, useGoTool bool) (*pluginpb.CodeGeneratorResponse, error) {
	if stderr == nil {
		stderr = os.Stderr
	}

	// Marshal the CodeGeneratorRequest.
	b, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}

	// Call the plugin with the modified CodeGeneratorRequest.
	cmdName := "protoc-gen-" + plugin
	var args []string

	if useGoTool {
		args = []string{"tool", cmdName}
		cmdName = "go"
	}

	var buf bytes.Buffer
	cmd := exec.Command(cmdName, args...)
	cmd.Stdin = bytes.NewReader(b)
	cmd.Stdout = &buf
	cmd.Stderr = stderr
	err = cmd.Run()
	if err != nil {
		return nil, err
	}

	// Read the CodeGeneratorResponse.
	var res pluginpb.CodeGeneratorResponse
	err = proto.Unmarshal(buf.Bytes(), &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// ReadRequest reads and unmarshals a CodeGeneratorRequest.
func ReadRequest(r io.Reader) (*pluginpb.CodeGeneratorRequest, error) {
	in, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	req := &pluginpb.CodeGeneratorRequest{}
	err = proto.Unmarshal(in, req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

// WriteResponse marshals and writes CodeGeneratorResponse res to w.
func WriteResponse(w io.Writer, res *pluginpb.CodeGeneratorResponse) error {
	out, err := proto.Marshal(res)
	if err != nil {
		return err
	}
	_, err = w.Write(out)
	return err
}
