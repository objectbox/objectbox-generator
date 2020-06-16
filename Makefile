# Default target executed when no arguments are given to make.
default_target: all

.PHONY: default_target clean depend build help

help:			## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'



#==============================================================================

all: depend build test

build:			## Build all targets
	CGO_ENABLED=1 go build ./cmd/objectbox-generator/

test: build		## Test all targets
	# echo "NOTE tests are WIP, currently not executed"
	# go test -v ./...

clean:			## Clean previous builds
	go clean -cache ./..
	rm -f objectbox-generator
	rm -f objectbox-generator.exe
	rm -rf third_party/flatbuffers-c-bridge/cmake-build

depend:			## Build dependencies
	./third_party/flatbuffers-c-bridge/build.sh
