BINARY  := mirakuta
PKG     := github.com/mirakuta-dev/mirakuta
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -s -w -X $(PKG)/internal/cli.Version=$(VERSION)

.PHONY: build test vet fmt clean cross dist

build:
	go build -ldflags "$(LDFLAGS)" -o $(BINARY).exe ./cmd/mirakuta

test:
	go test ./...

vet:
	go vet ./...

fmt:
	gofmt -w .

clean:
	rm -rf $(BINARY).exe dist/

cross: dist
	GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/$(BINARY)-windows-amd64.exe ./cmd/mirakuta
	GOOS=windows GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o dist/$(BINARY)-windows-arm64.exe ./cmd/mirakuta

dist:
	mkdir -p dist
