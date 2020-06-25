.DEFAULT: install

install:
	go install ./cmd/protoc-gen-go-patch

tools: internal/tools/*.go
	go generate --tags tools ./internal/tools

generate:
	go generate ./...

vet:
	go vet ./...

test: test-go test-cgo-disabled

test-go:
	go test -i -mod=readonly -race ./...
	go test -mod=readonly -v -race ./...

test-cgo-disabled:
	CGO_ENABLED=0 go test -i -mod=readonly ./...
	CGO_ENABLED=0 go test -mod=readonly -v ./...

go_module = $(shell go list -m)
proto_files = $(sort $(shell find . -name '*.proto'))

protos: $(proto_files)

.PHONY: $(proto_files)
$(proto_files): tools Makefile
	protoc \
		-I . \
		-I `go list -m -f {{.Dir}} google.golang.org/protobuf` \
		--go-patch_out=plugin=go,paths=import,module=$(go_module):. \
		--go-patch_out=plugin=go-grpc,paths=import,module=$(go_module):. \
		$@

.PHONY: shims
shims: shims/gogoproto

gogo_dir = $(shell go list -m -f {{.Dir}} github.com/gogo/protobuf)

.PHONY: shims/gogoproto
shims/gogoproto: tools Makefile
	protoc \
		-I . \
		-I $(gogo_dir) \
		--go_out=paths=source_relative:shims \
		$(gogo_dir)/gogoproto/gogo.proto
