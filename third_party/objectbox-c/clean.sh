#!/usr/bin/env bash
set -euo pipefail

scriptDir=$(dirname "${BASH_SOURCE[0]}")
set +x
rm -rf "${scriptDir}/download"
rm -rf "${scriptDir}/include"
rm -rf "${scriptDir}/lib"