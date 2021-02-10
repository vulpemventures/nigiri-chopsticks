.PHONY: build clean fmt run test

## build
build: clean 
	@echo "Build..."
	@go build -o build/nigiri-chopsticks

## remove the compiled binaries in ./bin
clean: 
	@echo "Clean..."
	@rm -rf ./build
	@rm -rf ./registry

## gofmt: Go Format
fmt:
	@echo "Gofmt..."
	@if [ -n "$(gofmt -l .)" ]; then echo "Go code is not formatted"; exit 1; fi

## help: prints this help message
help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

## run the short tests
test: fmt 
	@echo "Test..."
	CI=local go test -v -race ./...

