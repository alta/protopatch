package patch

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"strings"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
)

// StripParam strips a named param from req.
func StripParam(req *pluginpb.CodeGeneratorRequest, p string) {
	if req.Parameter == nil {
		return
	}
	v := stripParam(*req.Parameter, p)
	req.Parameter = &v

}

func stripParam(s, p string) string {
	var b strings.Builder
	for _, param := range strings.Split(s, ",") {
		if strings.SplitN(param, "=", 2)[0] != p {
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
func RunPlugin(plugin string, req *pluginpb.CodeGeneratorRequest, stderr io.Writer) (*pluginpb.CodeGeneratorResponse, error) {
	if stderr == nil {
		stderr = os.Stderr
	}

	// Marshal the CodeGeneratorRequest.
	b, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}

	// Call the plugin with the modified CodeGeneratorRequest.
	var buf bytes.Buffer
	cmd := exec.Command("protoc-gen-" + plugin)
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

// Write marshals and writes CodeGeneratorResponse res to w.
func Write(res *pluginpb.CodeGeneratorResponse, w io.Writer) error {
	out, err := proto.Marshal(res)
	if err != nil {
		return err
	}
	_, err = w.Write(out)
	return err
}
