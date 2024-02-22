#!/usr/bin/env python3

# Front-end Python script for managing docker image builds.
# 
# Images are 
#   .. defined as ci/docker/Dockerfile.<image-name>
#   .. built via `docker-build.py -i <image-name>`.
#   .. published using `-p` flag.
#

from os import system, listdir
from fnmatch import fnmatch
from os.path import realpath, dirname, abspath, relpath
from datetime import date
import glob
import argparse

topdir=relpath(dirname(realpath(__file__))+"/..")
docker_dir=topdir+"/ci/docker"

# Constants:
group="buildenv-generator"

# Defaults:
version=date.today()
image="ubuntu"

# CLI:
parser = argparse.ArgumentParser(
                    prog='docker-build',
                    description='Docker build shell',
                    epilog=f'See {docker_dir}')
parser.add_argument('-n','--dry-run', help='do not run', action='store_true')
parser.add_argument('-p','--push', help='push to repository', action='store_true')
parser.add_argument('--version', help=f"specify version (default: {version})", default=version)
parser.add_argument('-l','--list', help='list available docker images', action='store_true')
parser.add_argument('-i','--image', help=f"select image (default: {image})", default=image)
args = parser.parse_args()

# Dry-run helper:
def exec(cmd):
    if args.dry_run:
        print("DRY: "+cmd)
    else:
        system(cmd)

# Action: List
if args.list:
    images = [file[11:] for file in listdir(docker_dir) if fnmatch(file, 'Dockerfile.*')]
    print(f"Available images: {','.join(images)}")
    exit(0)

# Evaluate full image
full=f"objectboxio/{group}-{args.image}"

print(f"Building image: {full}")

# Build
exec(f"docker build -f {docker_dir}/Dockerfile.{args.image} -t {full}:{version} -t {full}:latest {topdir}")

# Push
if args.push:
    exec(f"docker push {full}:{version}")
    exec(f"docker push {full}:latest")
