#!/usr/bin/env bash
set -euo pipefail

# Usage: ./docker-build.sh [--version YYYY-MM-DD] [--push] [--push-existing] [--run]
# Note: order of arguments matters
# Builds the Docker image from Dockerfile.ubuntu in this directory.

# Parse args
version=$(date +%Y-%m-%d)
do_push=false
do_build=true
do_run=false

if [ "${1:-}" == "--version" ]; then
  version=$2
  shift 2
fi

if [ "${1:-}" == "--push" ]; then
  do_push=true
  shift
fi

if [ "${1:-}" == "--push-existing" ]; then
  do_push=true
  do_build=false
  shift
fi

# Start an interactive bash session in the image
if [ "${1:-}" == "--run" ]; then
  do_run=true
  shift
fi

# Image name
repository=objectboxio/buildenv-generator-ubuntu
full_name=$repository:$version

echo "------------------------------------------------------------------------"
echo "   Image name: $full_name"
echo "      Version: $version"
echo " Dockerfile: Dockerfile.ubuntu"
echo "------------------------------------------------------------------------"

if [ "$do_build" = true ]; then
  echo "Building \"$full_name\"..."
  docker build --tag "$full_name" \
    --build-arg version="$version" \
    --progress=plain \
    -f Dockerfile.ubuntu \
    .
  echo "Built \"$full_name\""
  if [ -f "./print-versions.sh" ]; then
    docker run --rm "$full_name" /print-versions.sh || true
  else
    docker run --rm "$full_name" cat /versions.txt || true
  fi
  docker image ls "$full_name"
fi

# Run an interactive shell in the image if requested
if [ "$do_run" = true ]; then
  echo "Starting interactive bash in \"$full_name\" (temporary container)..."
  docker run --rm -it "$full_name" bash
fi

if [ "$do_push" = true ]; then
  echo "Pushing \"$full_name\"..."
  docker push "$full_name"
  echo "Pushed \"$full_name\""
fi

# Print Docker version so we can copy it for documentation
docker --version
