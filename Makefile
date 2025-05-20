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


## run: Run a bitcoin version against local nigiri (docker stop chopsticks before running this)
run: clean build
	./build/nigiri-chopsticks --rpc-addr 127.0.0.1:18443 --electrs-addr localhost:30000 --use-faucet --use-mining --use-logger --addr localhost:3000


## runliquid: Run a liquid version against local nigiri (docker stop chopsticks-liquid before running this)
runliquid: clean build
	./build/nigiri-chopsticks --chain liquid --rpc-addr localhost:18884 --electrs-addr localhost:30001 --use-faucet --use-mining --use-logger --addr localhost:3001

## run the short tests
test: fmt 
	@echo "Test..."
	CI=false go test -v -race ./...

## run the CI short tests
testci: fmt 
	@echo "Test..."
	CI=true go test -v -race ./...
