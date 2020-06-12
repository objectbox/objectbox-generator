#!/usr/bin/env bash
set -euo pipefail

srcDir=$(dirname "${BASH_SOURCE[0]}")
buildDir=${srcDir}/cmake-build
installDir=${1:-${buildDir}}

function build() {
    echo "******** Configuring ********"
    cmake -S "$srcDir" -B "$buildDir" -DCMAKE_BUILD_TYPE=Release

    echo "******** Building ********"
    cmake --build "$buildDir" \
        --target flatbuffers \
        --target flatbuffers-c-bridge \
        --target flatbuffers-c-bridge-flatc \
        --target flatbuffers-c-bridge-test \
        -- -j
}

function install() {
    echo "******** Collecting artifacts ********"
    if [[ "${installDir}" != "${buildDir}" ]]; then
        echo "Copying from $buildDir to $installDir:"
    fi
    rsync -v \
        "$buildDir"/libflatbuffers-c-bridge.a \
        "$buildDir"/libflatbuffers-c-bridge-flatc.a \
        "$buildDir"/_deps/flatbuffers-*-build/libflatbuffers.a \
        "$installDir"
    echo "The compiled libraries can be found here:"
    ls -alh "$buildDir"/libflatbuffers*.a
}

function test() {
    echo "******** Testing ********"
    (cd "$buildDir" && ./flatbuffers-c-bridge-test)
}

build
test
install
