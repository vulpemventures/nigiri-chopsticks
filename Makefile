.PHONY: build clean fmt test

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

## run: Run a liquid version against local nigiri (docker stop chopsticks-liquid before running this)
run: clean build
	./build/nigiri-chopsticks --chain liquid --rpc-addr localhost:7041 --electrs-addr localhost:3012 --use-faucet --use-mining --use-logger --addr localhost:3001

## run the short tests
test: fmt 
	@echo "Test..."
	CI=local go test -v -race ./...

