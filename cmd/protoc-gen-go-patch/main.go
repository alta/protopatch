package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/alta/protopatch/patch"

	"google.golang.org/protobuf/compiler/protogen"
)

func main() {
	log.SetFlags(0)

	if len(os.Args) == 2 && os.Args[1] == "--version" {
		log.Printf("%v %v\n", filepath.Base(os.Args[0]), patch.Version)
		os.Exit(1)
		return
	}

	err := run()
	if err != nil {
		log.Printf("Error: %s", err)
		os.Exit(1)
	}
}

func run() error {
	req, err := patch.Read(os.Stdin)
	if err != nil {
		return err
	}

	var plugin string

	opts := protogen.Options{
		ParamFunc: func(name, value string) error {
			switch name {
			case "plugin":
				plugin = value
			}
			return nil // Ignore unknown params.
		},
	}

	gen, err := opts.New(req)
	if err != nil {
		return err
	}

	if plugin == "" {
		s := strings.TrimPrefix(filepath.Base(os.Args[0]), "protoc-gen-")
		return fmt.Errorf("no protoc plugin specified; use 'protoc --%s_out=plugin=$PLUGIN:...'", s)
	}

	if os.Getenv("PROTO_PATCH_DEBUG_LOGGING") == "" {
		log.SetOutput(ioutil.Discard)
	}

	// Strip our custom param(s).
	patch.StripParam(gen.Request, "plugin")

	// Run the specified plugin and unmarshal the CodeGeneratorResponse.
	res, err := patch.RunPlugin(plugin, gen.Request, nil)
	if err != nil {
		return err
	}

	// Initialize a Patcher and scan source proto files.
	patcher, err := patch.NewPatcher(gen)
	if err != nil {
		return err
	}

	// Patch the CodeGeneratorResponse.
	err = patcher.Patch(res)
	if err != nil {
		return err
	}

	log.Printf("Writing %d file(s) to protoc: %s", len(res.File), *res.File[0].Name)

	// Write the patched CodeGeneratorResponse to stdout.
	return patch.Write(res, os.Stdout)
}
