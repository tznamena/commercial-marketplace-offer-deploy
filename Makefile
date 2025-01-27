# generate runs `go generate` to build the dynamically generated
# source files, except the protobuf stubs which are built instead with
# "make protobuf".
generate:
	./scripts/generate-code.sh

clean:
	rm -rf bin
	mkdir -p bin

build:
	go build -o ./bin/ ./cmd/apiserver
	go build -o ./bin/ ./cmd/operator
	make tools

apiserver-local:
	./scripts/run-local.sh apiserver

operator-local:
	./scripts/run-local.sh operator

# Builds docker container, starts ngrok in the background, and 
# calls docker compose up with the public NGROK endpoint for MODM to receive event messages from Azure
run-local:
	./scripts/run-local.sh modm $(build)

run-testharness:
	make run -C ./tools

test:
	go test ./...

test-integration:
	$(ENV_LOCAL_TEST) \
	go test -tags=integration ./test -v -count=1 

sdk:
	go build ./sdk

tools:
	./scripts/build-tools.sh

assemble: apiserver operator 
	./scripts/assemble.sh

.NOTPARALLEL:

.PHONY: build test run run-testharness apiserver-local operator-local sdk generate tools
