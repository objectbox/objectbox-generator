#!/usr/bin/env bash
# Script to update ObjectBox C version across the project
set -euo pipefail

# macOS does not have realpath and readlink does not have -f option, so do this instead:
script_dir=$( cd "$(dirname "$0")" ; pwd -P )
src_dir="$script_dir/.."

echo "set-objectbox-c-version.sh script dir: $script_dir"

if [[ "$#" -ne "1" ]]; then
  echo "usage: $0 <version>"
  echo "e.g. $0 4.3.1"
  echo "  or $0 5.0.0-rc"
  exit 1
fi

version=$1

echo "Setting ObjectBox C version: $version"

# align GNU vs BSD `sed` version handling -i argument
if [[ "$OSTYPE" == "darwin"* ]]; then
  sed="sed -i ''"
else
  sed="sed -i"
fi

function replace() {
  file=$1
  pattern=$2

  # `sed` doesn't fail if the file isn't changed - verify the checksum instead
  echo "Updating file $file with pattern $pattern"
  checksum=$(sha1sum "$file" 2>/dev/null || shasum -a 1 "$file")
  $sed "$pattern" "$file"

  new_checksum=$(sha1sum "$file" 2>/dev/null || shasum -a 1 "$file")
  if [[ "$checksum" == "$new_checksum" ]]; then
    echo "No change to $file - was it already at the intended version? (Otherwise check manually)"
  fi
  echo ""
}

# Update CMake files with GIT_TAG
replace "$src_dir/examples/cpp/CMakeLists.txt" "s|GIT_TAG[[:space:]]*v[0-9][0-9.rc-]*|GIT_TAG        v${version}|g"
replace "$src_dir/examples/cpp-flat-layout/CMakeLists.txt" "s|GIT_TAG[[:space:]]*v[0-9][0-9.rc-]*|GIT_TAG        v${version}|g"
replace "$src_dir/test/integration/cmake/projects/common.cmake" "s|GIT_TAG[[:space:]]*v[0-9][0-9.rc-]*|GIT_TAG        v${version}|g"

# Update bash script with cVersion variable
replace "$src_dir/third_party/objectbox-c/get-objectbox-c.sh" "s/^cVersion=.*/cVersion=${version}/g"

echo "Version update complete!"
