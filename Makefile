BINARY := podgo

.PHONY: all build install run test test-cover fmt vet lint tidy clean deps

all: build

## build: compile the binary
build:
	go build -o $(BINARY) .

## install: install the binary into $GOBIN (or $GOPATH/bin)
install:
	go install .

## run: build and run the app
run: build
	./$(BINARY)

## test: run the test suite
test:
	go test ./...

## test-cover: run tests with a coverage report
test-cover:
	go test -cover ./...

## fmt: format the code
fmt:
	gofmt -w .

## vet: run go vet
vet:
	go vet ./...

## lint: run golangci-lint (requires golangci-lint installed)
lint:
	golangci-lint run

## tidy: sync go.mod / go.sum
tidy:
	go mod tidy

## deps: system dependencies required to build (Debian/Ubuntu)
deps:
	sudo apt-get update && sudo apt-get install -y libasound2-dev

## clean: remove build artifacts
clean:
	rm -f $(BINARY)
