#!/usr/bin/env bash
set -euo pipefail

#fccVersion=v0.6.2
fccVersion=f064cefb2034d1e7407407ce32a6085c322212a7 # We need CMake 4.x support, for which no release exists yet
fccRepo=github.com/dvidelabs/flatcc

scriptDir=$(dirname "${BASH_SOURCE[0]}")
srcDir=${scriptDir}/src
buildDir=${scriptDir}/build
installDir=${scriptDir}

buildType=Release
configArgs="-DCMAKE_BUILD_TYPE=${buildType}"

if [[ "$(uname)" == MINGW* ]] || [[ "$(uname)" == CYGWIN* ]]; then
    configArgs+=' -G "MinGW Makefiles"'
    # aligned_alloc() would be missing on MinGW, see https://github.com/dvidelabs/flatcc/issues/155
    # It's fixed now on master, the following line can be removed after upgrading to a new release
    export CFLAGS="-DFLATCC_USE_GENERIC_ALIGNED_ALLOC=1"
fi

function prepare() {
    echo "******** Getting Flatcc sources ********"
    if [[ ! -d ${srcDir} ]]; then
        echo "Cloning ${fccRepo} into Into: ${srcDir}"
        git clone https://${fccRepo}.git "${srcDir}"
    fi

    echo "Checking out ${fccVersion}"
    (cd "${srcDir}"; git fetch)
    (cd "${srcDir}"; git checkout ${fccVersion})
}

function build() {
    echo "******** Configuring & building ********"

    # Note: we need an absolute path...
    # realpath isn't available on macOS and the "else" variant didn't work well on windows because the path was already absolute...
    srcDirAbsolute=
    if [[ -x $(command -v realpath) ]]; then
        srcDirAbsolute=$(realpath "$srcDir")
    else
        srcDirAbsolute="$(pwd)/$srcDir"
    fi

    pwd=$(pwd)
    mkdir -p "$buildDir"

    set -x

    # need to use eval because of quotes in configArgs... bash is just wonderful...
    cd "$buildDir"
    eval "cmake \"$srcDirAbsolute\" $configArgs"
    cd "$pwd"

    cmake --build "$buildDir" --config ${buildType} --target flatccrt
    set +x
}

function install() {
    echo "******** Collecting artifacts ********"
    echo "Copying from ${srcDir} to ${installDir}:"
    cp -rv "${srcDir}/include" "${installDir}"
    if [[ -d "${srcDir}/lib/${buildType}" ]]; then
        cp -rv "${srcDir}/lib/${buildType}" "${installDir}/lib"
    else
        cp -rv "${srcDir}/lib" "${installDir}"
    fi
}

prepare
build
install
