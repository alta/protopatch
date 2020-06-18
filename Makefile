.DEFAULT: install

install:
	go install ./cmd/protoc-gen-go-patch

test:
	go test -race ./...
	CGO_ENABLED=0 go test ./...

tools: internal/tools/*.go
	go generate --tags tools ./internal/tools

protos: tools
	protoc \
		-I . \
		-I `go list -m -f {{.Dir}} google.golang.org/protobuf` \
		--go-patch_out=plugin=go,paths=source_relative:. \
		--go-patch_out=plugin=go-grpc,paths=source_relative:. \
		patch/*.proto
