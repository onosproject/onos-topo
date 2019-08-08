export CGO_ENABLED=0
export GO111MODULE=on

.PHONY: build

ONOS_TOPO_VERSION := latest
ONOS_TOPO_DEBUG_VERSION := debug
ONOS_BUILD_VERSION := stable

build: # @HELP build the Go binaries and run all validations (default)
build:
	CGO_ENABLED=1 go build -o build/_output/onos-topo ./cmd/onos-topo
	CGO_ENABLED=1 go build -gcflags "all=-N -l" -o build/_output/onos-topo-debug ./cmd/onos-topo

test: # @HELP run the unit tests and source code validation
test: build deps license_check linters
	go test github.com/onosproject/onos-topo/pkg/...
	go test github.com/onosproject/onos-topo/cmd/...

coverage: # @HELP generate unit test coverage data
coverage: build deps linters license_check
	./build/bin/coveralls-coverage

deps: # @HELP ensure that the required dependencies are in place
	go build -v ./...
	bash -c "diff -u <(echo -n) <(git diff go.mod)"
	bash -c "diff -u <(echo -n) <(git diff go.sum)"

linters: # @HELP examines Go source code and reports coding problems
	golangci-lint run

license_check: # @HELP examine and ensure license headers exist
	./build/licensing/boilerplate.py -v

protos: # @HELP compile the protobuf files (using protoc-go Docker)
	docker run -it -v `pwd`:/go/src/github.com/onosproject/onos-topo \
		-w /go/src/github.com/onosproject/onos-topo \
		--entrypoint build/bin/compile-protos.sh \
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
