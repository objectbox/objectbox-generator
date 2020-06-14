#!/usr/bin/env bash
set -euo pipefail

srcDir=$(dirname "${BASH_SOURCE[0]}")
buildDir=${srcDir}/cmake-build
installDir=${1:-${buildDir}}

buildType=Release

libPrefix=lib
libExt=.a
exeExt=
configArgs="-DCMAKE_BUILD_TYPE=${buildType}"
buildArgs="-- -j"
buildOutputDir=
if [[ "$(uname)" == MINGW* ]] || [[ "$(uname)" == CYGWIN* ]]; then
    exeExt=.exe

    # MinGW build
    configArgs+=' -G "MinGW Makefiles"'

    # MSVC build - CGO linking fails with "undefined reference to '__CxxFrameHandler4'"
    #  libExt=.lib
    #  configArgs="-A x64"
    #  buildOutputDir=/${buildType}
    #  # buildArgs="-- /m"    fails with "error MSB1008: Only one project can be specified."
fi

function build() {
    echo "******** Configuring & building ********"
    set -x
    # need to use eval because of quotes in configArgs... bash is just wonderful...
    eval "cmake -S \"$srcDir\" -B \"$buildDir\" $configArgs"

    # Note: flatbuffers-c-bridge-test implies flatbuffers, flatbuffers-c-bridge and flatbuffers-c-bridge-flatc
    # We don't specify them explicitly to be compatible with MSVC which allows only one target per cmake call...
    cmake --build "$buildDir" --config ${buildType} --target flatbuffers-c-bridge-test ${buildArgs}
    set +x
}

function install() {
    echo "******** Collecting artifacts ********"
    if [[ "${installDir}" != "${buildDir}${buildOutputDir}" ]]; then
        echo "Copying from ${buildDir}${buildOutputDir} to ${installDir}:"
        cp "${buildDir}${buildOutputDir}"/${libPrefix}flatbuffers-c-bridge${libExt} "$installDir"
        cp "${buildDir}${buildOutputDir}"/${libPrefix}flatbuffers-c-bridge-flatc${libExt} "$installDir"
    fi
    cp "${buildDir}"/_deps/flatbuffers-*-build${buildOutputDir}/${libPrefix}flatbuffers${libExt} "$installDir"
    echo "The compiled libraries can be found here:"
    ls -alh "$installDir"/${libPrefix}flatbuffers*${libExt}
}

function test() {
    echo "******** Testing ********"
    (cd "${buildDir}${buildOutputDir}" && ./flatbuffers-c-bridge-test${exeExt})
}

build
test
install
