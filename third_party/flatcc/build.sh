#!/usr/bin/env bash
set -euo pipefail

fccVersion=v0.6.0
fccRepo=github.com/dvidelabs/flatcc

scriptDir=$(dirname "${BASH_SOURCE[0]}")
srcDir=${scriptDir}/src
buildDir=${scriptDir}/build
installDir=${scriptDir}

buildType=Release
configArgs="-DCMAKE_BUILD_TYPE=${buildType}"

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

    srcDirAbsolute="$(pwd)/$srcDir"
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
