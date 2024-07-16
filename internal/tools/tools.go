//go:build tools

package tools

import (
	_ "github.com/envoyproxy/protoc-gen-validate"
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
)

//go:generate go install github.com/envoyproxy/protoc-gen-validate
//go:generate go install google.golang.org/protobuf/cmd/protoc-gen-go
//go:generate go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
//go:generate go install ../../cmd/protoc-gen-go-patch
