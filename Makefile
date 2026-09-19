BINARY := micoprompt
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -s -w -X main.version=$(VERSION)

.PHONY: build run test lint clean release

build:
	go build -ldflags "$(LDFLAGS)" -o bin/$(BINARY) ./cmd

run: build
	./bin/$(BINARY) -dir . -out micoprompt.txt

test:
	go test ./...

lint:
	golangci-lint run

clean:
	rm -rf bin/ dist/ prompt.txt

# кросс-компиляция
release:
	GOOS=linux   GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/$(BINARY)-linux-amd64    ./cmd/
	GOOS=darwin  GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o dist/$(BINARY)-darwin-arm64   ./cmd/
	GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/$(BINARY)-windows-amd64.exe ./cmd/