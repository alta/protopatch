.DEFAULT: install

.PHONY: install tools generate vet test test-go test-cgo-disabled protos

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

proto_files = $(sort $(shell find patch -name '*.proto'))

protos: $(proto_files)

.PHONY: $(proto_files)
$(proto_files): Makefile
	protoc \
		-I . \
		-I `go list -m -f {{.Dir}} google.golang.org/protobuf` \
		--go-patch_out=plugin=go,paths=import,module=`go list -m`:. \
		--go-patch_out=plugin=go-grpc,paths=import,module=`go list -m`:. \
		$@
