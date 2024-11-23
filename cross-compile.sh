#!/bin/bash
# This script cross-compiles the do-ssh binary for multiple architectures and operating systems.

if git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
    VERSION=$(git describe --tags --always)
else
    echo "Error: This script must be run in a git repository."
    exit 1
fi
OSs="linux darwin"
VERSION=$(git describe --tags --always)

for GOOS in $OSs; do
    for GOARCH in $ARCHS; do
        docker run --rm -it \
            -v "$PWD":/usr/src/do-ssh \
            -e GOOS=$GOOS \
            -e GOARCH=$GOARCH \
            -e VERSION=$VERSION \
            -w /usr/src/do-ssh \
            golang:1.23 bash -c "go build -v -o do-ssh-$GOOS-$GOARCH-$VERSION"
    done
done
