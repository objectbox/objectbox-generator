# Default target executed when no arguments are given to make.
default_target: all

.PHONY: default_target help clean depend build test test-depend

help:			## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

#==============================================================================

all: depend build

build:			## Build all targets
	CGO_ENABLED=1 go build ./cmd/objectbox-generator/

test: 			## Test all targets
	go test ./...

clean:			## Clean previous builds
	go clean -cache ./..
	rm -f objectbox-generator
	rm -f objectbox-generator.exe
	rm -rf third_party/flatbuffers-c-bridge/cmake-build
	./third_party/flatcc/clean.sh
	./third_party/objectbox-c/clean.sh

depend:			## Build dependencies
	./third_party/flatbuffers-c-bridge/build.sh

test-depend: depend		## Build test dependencies
	./third_party/flatcc/build.sh
	./third_party/objectbox-c/get-objectbox-c.sh

