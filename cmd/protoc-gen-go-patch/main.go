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

	if os.Getenv("PROTO_PATCH_DEBUG_LOGGING") == ""  {
		log.SetOutput(ioutil.Discard)
	}

	var plugin string

	protogen.Options{
		ParamFunc: func(name, value string) error {
			switch name {
			case "plugin":
				plugin = value
			}
			return nil // Ignore unknown params.
		},
	}.Run(func(gen *protogen.Plugin) error {
		if plugin == "" {
			s := strings.TrimPrefix(filepath.Base(os.Args[0]), "protoc-gen-")
			return fmt.Errorf("no protoc plugin specified; use 'protoc --%s_out=plugin=$PLUGIN:...'", s)
		}

		// Strip our custom param(s).
		patch.StripParam(gen.Request, "plugin")

		// Run the specified plugin and unmarshal the CodeGeneratorResponse.
		res, err := patch.RunPlugin(plugin, gen.Request, nil)
		if err != nil {
			return err
		}

		// Create a cache
		cache := patch.NewCache()
		switch plugin {
		// Cache the protoc-gen-go generated base definitions for later use
		case "go":
			if err := cache.Save(res); err != nil {
				return err
			}
		// Other plugins needs the base definitions in order to apply the renaming
		default:
			if err := cache.Load(res); err != nil {
				return err
			}
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

		// Remove added cached files
		switch plugin {
		case "go":
		default:
			cache.CleanResFiles(res)
		}
		// Write the patched CodeGeneratorResponse to stdout.
		if err := patch.Write(res, os.Stdout); err != nil {
			return err
		}
		return nil
	})
}
