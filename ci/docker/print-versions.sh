#!/usr/bin/env bash
set -euo pipefail

echo "=== Environment ==="
uname -a || true
echo "PATH=$PATH"

echo
echo "=== Go versions (system installed first) ==="
echo "which: /usr/bin/go"
/usr/bin/go version || true

# Custom (manual) installs
for go_custom_install in /usr/local/go*/bin/go; do
  if [ -x "$go_custom_install" ]; then
    echo "which: $go_custom_install"
    "$go_custom_install" version || true
  fi
done

# Also list Go toolchains installed via `go install golang.org/dl/...` (e.g., /root/go/bin/go1.24.9)
for go_dl_tool in /root/go/bin/go1.*; do
  if [ -x "$go_dl_tool" ]; then
    echo "which: $go_dl_tool"
    "$go_dl_tool" version || true
  fi
done

echo
echo "=== Compiler, libs & Tools ==="
lib_dir="/usr/lib/x86_64-linux-gnu"
"$lib_dir/libc.so.6" | grep "GLIBC"
strings "$lib_dir/libstdc++.so.6" | grep GLIBCXX_3 | head -1
strings "$lib_dir/libstdc++.so.6" | grep GLIBCXX_3 | tail -1
which gcc
gcc --version | grep gcc
which ld
ld --version | grep "ld"
cmake --version | grep version
make --version | grep Make
ccache --version | grep "ccache version"
