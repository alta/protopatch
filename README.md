# protopatch

Utility to patch `protoc-gen-go` output with Go-specific features. WIP.

## Install

`go install github.com/alta/protopatch/cmd/protoc-gen-go-patch`

## Usage

After installing `protoc-gen-go-patch`, use it by specifying it with a `--go-patch_out=...` argument to `protoc`:

```shell
protoc \
	-I . \
	-I `go list -m -f {{.Dir}} google.golang.org/protobuf` \
	--go-patch_out=plugin=go,paths=source_relative:. \
	--go-patch_out=plugin=go-grpc,paths=source_relative:. \
	*.proto
```
