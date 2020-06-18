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

protos: tools
	protoc \
		-I . \
		-I `go list -m -f {{.Dir}} google.golang.org/protobuf` \
		--go-patch_out=plugin=go,paths=source_relative:. \
		--go-patch_out=plugin=go-grpc,paths=source_relative:. \
		patch/*.proto
