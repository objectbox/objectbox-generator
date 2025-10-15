#!/usr/bin/env bash
set -euo pipefail

cVersion=5.0.0-rc

scriptDir=$(dirname "${BASH_SOURCE[0]}")


echo "******** Downloading ObjectBox-C library ********"
echo "Into: ${scriptDir}"
cd "${scriptDir}"

# don't install the library system-wide, just download it
export installLibrary=false
bash <(curl -s https://raw.githubusercontent.com/objectbox/objectbox-c/main/download.sh) --quiet ${cVersion}

echo "******** Collecting artifacts ********"
cp -rfv download/*objectbox*/include ./
cp -rfv download/*objectbox*/lib ./
rm -rfv download

echo "Downloaded ObjectBox-C headers and library:"
echo "Current directory: ${scriptDir}"
ls -lh include lib