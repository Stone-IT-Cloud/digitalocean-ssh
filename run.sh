#!/bin/bash
# This script sets the basic variables to run a golang code in a docker container.

cmd="$*"

docker run --rm -it \
    -v "$PWD":/usr/src/docmd \
    -w /usr/src/docmd \
    golang:1.23 bash -c "$cmd"
