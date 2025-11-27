#!/usr/bin/env bash
set -euo pipefail

# This file is used to regenerate FlatBuffers reflection Go sources in our flatbuffersc Go wrapper

scriptDir=$(dirname "${BASH_SOURCE[0]}")
repoDir=$(realpath "${scriptDir}/../")
tmpDir="${repoDir}"/scripts/tmp-generated

# Note: we're using objectbox-generator's flatc-integration.
# If you can't build objectbox-generator, you can build and use flatbuffers flatc binary instead.
flatc="go run ./cmd/objectbox-generator FLATC"
# flatc="../objectbox-generator FLATC" # You can use this if objectbox-generator was already build locally

if [[ "${repoDir}" == "" ]]; then
    echo "Invalid repo directory detected"
    exit 1
fi

printf "\n******** Preparing temporary dir ********\n"
if [[ -d "${tmpDir}" ]]; then
    rm -rv "${tmpDir}"
fi

mkdir -pv "${tmpDir}"

printf "\n******** Generating Go code ********\n"
${flatc} --go -o "${tmpDir}" "${repoDir}"/third_party/flatbuffers-c-bridge/third_party/flatbuffers/reflection/reflection.fbs
echo "Generated files:"
ls "${tmpDir}/reflection"

printf "\n******** Updating repo code ********\n"
rsync -av --delete \
    "${tmpDir}/reflection/" \
    "${repoDir}/internal/generator/flatbuffersc/reflection/"

printf "\n******** Removing temporary dir ********\n"
rm -rv "${tmpDir}"