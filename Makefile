export CGO_ENABLED=0
export GO111MODULE=on

.PHONY: build

ONOS_TOPO_VERSION := latest
ONOS_TOPO_DEBUG_VERSION := debug
ONOS_BUILD_VERSION := stable

build: # @HELP build the Go binaries and run all validations (default)
build:
	CGO_ENABLED=1 go build -o build/_output/onos-topo ./cmd/onos
	CGO_ENABLED=1 go build -gcflags "all=-N -l" -o build/_output/onos-topo-debug ./cmd/onos-topo
	go build -o build/_output/onos ./cmd/onos

test: # @HELP run the unit tests and source code validation
test: build deps lint vet license_check gofmt cyclo misspell ineffassign
	go test github.com/onosproject/onos-topo/pkg/...
	go test github.com/onosproject/onos-topo/cmd/...

coverage: # @HELP generate unit test coverage data
coverage: build deps lint vet license_check gofmt cyclo misspell ineffassign
	./build/bin/coveralls-coverage

deps: # @HELP ensure that the required dependencies are in place
	go build -v ./...
	bash -c "diff -u <(echo -n) <(git diff go.mod)"
	bash -c "diff -u <(echo -n) <(git diff go.sum)"

lint: # @HELP run the linters for Go source code
	golint -set_exit_status github.com/onosproject/onos-topo/pkg/...
	golint -set_exit_status github.com/onosproject/onos-topo/cmd/...
	golint -set_exit_status github.com/onosproject/onos-topo/test/...

vet: # @HELP examines Go source code and reports suspicious constructs
	go vet github.com/onosproject/onos-topo/pkg/...
	go vet github.com/onosproject/onos-topo/cmd/...
	go vet github.com/onosproject/onos-topo/test/...

cyclo: # @HELP examines Go source code and reports complex cycles in code
	gocyclo -over 25 pkg/
	gocyclo -over 25 cmd/
	gocyclo -over 25 test/

misspell: # @HELP examines Go source code and reports misspelled words
	misspell -error -source=text pkg/
	misspell -error -source=text cmd/
	misspell -error -source=text test/
	misspell -error docs/

ineffassign: # @HELP examines Go source code and reports inefficient assignments
	ineffassign pkg/
	ineffassign cmd/
	ineffassign test/

license_check: # @HELP examine and ensure license headers exist
	./build/licensing/boilerplate.py -v

gofmt: # @HELP run the Go format validation
	bash -c "diff -u <(echo -n) <(gofmt -d pkg/ cmd/ tests/)"

protos: # @HELP compile the protobuf files (using protoc-go Docker)
	docker run -it -v `pwd`:/go/src/github.com/onosproject/onos-topo \
		-w /go/src/github.com/onosproject/onos-topo \
		--entrypoint pkg/northbound/proto/compile-protos.sh \
		onosproject/protoc-go:stable

onos-topo-base-docker: # @HELP build onos-topo base Docker image
	@go mod vendor
	docker build . -f build/base/Dockerfile \
		--build-arg ONOS_BUILD_VERSION=${ONOS_BUILD_VERSION} \
		-t onosproject/onos-topo-base:${ONOS_TOPO_VERSION}
	@rm -rf vendor

onos-topo-docker: onos-topo-base-docker # @HELP build onos-topo Docker image
	docker build . -f build/onos-topo/Dockerfile \
		--build-arg ONOS_TOPO_BASE_VERSION=${ONOS_TOPO_VERSION} \
		-t onosproject/onos-topo:${ONOS_TOPO_VERSION}

onos-topo-debug-docker: onos-topo-base-docker # @HELP build onos-topo Docker debug image
	docker build . -f build/onos-topo-debug/Dockerfile \
		--build-arg ONOS_TOPO_BASE_VERSION=${ONOS_TOPO_VERSION} \
		-t onosproject/onos-topo:${ONOS_TOPO_DEBUG_VERSION}

onos-cli-docker: onos-topo-base-docker # @HELP build onos-cli Docker image
	docker build . -f build/onos-cli/Dockerfile \
		--build-arg ONOS_TOPO_BASE_VERSION=${ONOS_TOPO_VERSION} \
		-t onosproject/onos-cli:${ONOS_TOPO_VERSION}

onos-topo-it-docker: onos-topo-base-docker # @HELP build onos-topo-integration-tests Docker image
	docker build . -f build/onos-it/Dockerfile \
		--build-arg ONOS_TOPO_BASE_VERSION=${ONOS_TOPO_VERSION} \
		-t onosproject/onos-topo-integration-tests:${ONOS_TOPO_VERSION}

# integration: @HELP build and run integration tests
integration: kind
	onit create cluster
	onit add simulator
	onit add simulator
	onit run suite integration-tests


images: # @HELP build all Docker images
images: build onos-topo-docker onos-topo-debug-docker

all: build images


clean: # @HELP remove all the build artifacts
	rm -rf ./build/_output ./vendor ./cmd/onos-topo/onos-topo ./cmd/dummy/dummy

help:
	@grep -E '^.*: *# *@HELP' $(MAKEFILE_LIST) \
    | sort \
    | awk ' \
        BEGIN {FS = ": *# *@HELP"}; \
        {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}; \
    '
