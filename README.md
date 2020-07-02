# protopatch

Patch `protoc` plugin output with Go-specific features. The `protoc-gen-go-patch` command wraps calls to Go code generators like [`protoc-gen-go`](https://pkg.go.dev/google.golang.org/protobuf/cmd/protoc-gen-go) or [`protoc-gen-go-grpc`](https://pkg.go.dev/google.golang.org/grpc/cmd/protoc-gen-go-grpc) and patches the Go syntax before being written to disk.

## Features

Patches are defined via an `Options` extension on messages, fields, `oneof` fields, enums, and enum values.

- `go.message` — message options, which modify the generated Go struct for a message.
- `go.field` — message field options, which modify Go struct fields and getter methods.
- `go.oneof` — `oneof` field options, which modify struct fields, interface types, and wrapper types.
- `go.enum` — `enum` options, which modify Go enum types and values.
- `go.value` — `enum` value options, which modify Go `const` values.

### Custom Names

```proto
import "patch/go.proto";

message OldName {
	option (go.message) = {name: 'NewName'};
	int id = 1 [(go.field) = {name: 'ID'}];
}

enum Errors {
	option (go.enum) = {name: 'ProtocolErrors'};
	INVALID = 1 [(go.value) = {name: 'ErrInvalid'}];
	NOT_FOUND = 2 [(go.value) = {name: 'ErrNotFound'}];
	TOO_FUN = 3 [(go.value) = {name: 'ErrTooFun'}];
}
```

### Struct Tags

```proto
message ToDo {
	int32 id = 1 [(go.field) = {name: 'ID', tags: '`xml:"id,attr"`'}];
	string description = 2 [(go.field) = {tags: '`xml:"desc"`'}];
}
```

### TODO: More Documentation

…

## Install

`go install github.com/alta/protopatch/cmd/protoc-gen-go-patch`

## Usage

After installing `protoc-gen-go-patch`, use it by specifying it with a `--go-patch_out=...` argument to `protoc`:

```shell
protoc \
	-I . \
	-I `go list -m -f {{.Dir}} github.com/alta/protopatch` \
	-I `go list -m -f {{.Dir}} google.golang.org/protobuf` \
	--go-patch_out=plugin=go,paths=source_relative:. \
	--go-patch_out=plugin=go-grpc,paths=source_relative:. \
	*.proto
```
