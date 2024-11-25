#!/bin/bash
# This script sets the basic variables to run a golang code in a docker container.

docker run --rm -it \
    -v "$PWD":/usr/src/docmd \
    -w /usr/src/docmd \
    golang:1.23 bash -c "${1}"
