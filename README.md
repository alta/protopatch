# protopatch

[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white)](https://pkg.go.dev/github.com/alta/protopatch) [![build status](https://img.shields.io/github/workflow/status/alta/protopatch/go.yaml.svg?branch=main)](https://github.com/alta/protopatch/actions)

Patch `protoc` plugin output with Go-specific features. The `protoc-gen-go-patch` command wraps calls to Go code generators like [`protoc-gen-go`](https://pkg.go.dev/google.golang.org/protobuf/cmd/protoc-gen-go) or [`protoc-gen-go-grpc`](https://pkg.go.dev/google.golang.org/grpc/cmd/protoc-gen-go-grpc) and patches the Go syntax before being written to disk.

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

## Features

Patches are defined via an `Options` extension on messages, fields, `oneof` fields, enums, and enum values.

- `go.message` — message options, which modify the generated Go struct for a message.
- `go.field` — message field options, which modify Go struct fields and getter methods.
- `go.oneof` — oneof field options, which modify struct fields, interface types, and wrapper types.
- `go.enum` — enum options, which modify Go enum types and values.
- `go.value` — enum value options, which modify Go const values.

### Custom Names

```proto
import "patch/go.proto";

message OldName {
	option (go.message).name = 'NewName';
	int id = 1 [(go.field).name = 'ID'];
}

enum Error {
	option (go.enum).name = 'ProtocolError';
	INVALID = 1 [(go.value).name = 'ErrInvalid'];
	NOT_FOUND = 2 [(go.value).name = 'ErrNotFound'];
	TOO_FUN = 3 [(go.value).name = 'ErrTooFun'];
}
```

#### Alternate Syntax

Multiple options can be grouped together with a message bounded by `{}`:

```proto
import "patch/go.proto";

message OldName {
	option (go.message) = {name: 'NewName'};
	int32 id = 1 [(go.field) = {name: 'ID'}];
}

enum Error {
	option (go.enum) = {name: 'ProtocolError'};
	INVALID = 1 [(go.value) = {name: 'ErrInvalid'}];
	NOT_FOUND = 2 [(go.value) = {name: 'ErrNotFound'}];
	TOO_FUN = 3 [(go.value) = {name: 'ErrTooFun'}];
}
```

### Struct Tags

```proto
message ToDo {
	int32 id = 1 [(go.field).name: 'ID', (go.field).tags = 'xml:"id,attr"'];
	string description = 2 [(go.field).tags: 'xml:"desc"'];
}
```

#### Alternate Syntax

Multiple options can be grouped together with a message bounded by `{}`:

```proto
message ToDo {
	int32 id = 1 [(go.field) = {name: 'ID', tags: 'xml:"id,attr"'}];
	string description = 2 [(go.field) = {tags: 'xml:"desc"'}];
}
```

### Embedded Fields

A message field can be embedded in the generated [Go struct](https://golang.org/ref/spec#Struct_types) with the `(go.field).embed` option. This only works for message fields, and will not work for oneof fields or basic types.

```proto
import "patch/go.proto";

message A {
	B b = 1 [(go.field).embed = true];
}

message B {
	string value = 1;
}
```

The resulting Go struct will partially have the form:

```go
type A struct {
	*B
}

type B struct {
	Value string
}

var a A
a.Value = "value" // This works because B is embedded in A
```

#### Alternate Syntax

Multiple options can be grouped together with a message bounded by `{}`:

```proto
import "patch/go.proto";

message A {
	B b = 1 [(go.field) = {embed: true}];
}

message B {
	string value = 1;
}
```

### Linting

Protopatch can automatically “lint” generated names into something resembling [idiomatic Go style](https://golang.org/doc/effective_go.html#names). This feature should be considered *unstable*, and the names it generates are subject to change as this feature evolves.

- **Initialisms:** names with `ID` or `URL` or other well-known initialisms will have their case preserved. For example `Id` would lint to `ID`, and `ApiBaseUrl` would lint to `APIBaseURL`.
- **Stuttering:** it will attempt to remove repeated prefixed names from enum values. An enum value of type `Foo` named `Foo_FOO_BAR` would lint to `FooBar`.

To lint all generated Go names, add `option (go.lint).all = true` to your `proto` file. To lint only enum values, add `option (go.lint).values = true`. To specify one or more custom initialisms, specify an initialism with `option (go.lint).initialisms = 'HSV'` for the `HSV` initialism. All names with `HSV` will preserve its case.

```proto
option (go.lint).all = true;
option (go.lint).initialisms = 'RGB';
option (go.lint).initialisms = 'RGBA';
option (go.lint).initialisms = 'HSV';

enum Protocol {
	// PROTOCOL_INVALID value should lint to ProtocolInvalid.
	PROTOCOL_INVALID = 0;
	// PROTOCOL_IP value should lint to ProtocolIP.
	PROTOCOL_IP = 1;
	// PROTOCOL_UDP value should lint to ProtocolUDP.
	PROTOCOL_UDP = 2;
	// PROTOCOL_TCP value should lint to ProtocolTCP.
	PROTOCOL_TCP = 3;
}

message Color {
	oneof value {
		// rgb should lint to RGB.
		string rgb = 1;
		// rgba should lint to RGBA.
		string rgba = 2;
		// hsv should lint to HSV.
		string hsv = 3;
	}
}
```
