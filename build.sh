#!/bin/bash
set -e
ROOT=$(cd "$(dirname "$0")"; pwd)
NAME=$(echo $TRAVIS_REPO_SLUG | sed 's|.*/||')
RELEASE=$(git describe --always --tags)

echo "Building ..."
mkdir build dist
gox -output "build/{{.OS}}_{{.Arch}}/{{.Dir}}"

echo "Packaging ..."
cd build
for OSARCH in $(ls); do
	cd "$OSARCH"
	tar -czf "$ROOT/dist/${NAME}_${RELEASE}_${OSARCH}.tar.gz" .
	cd -
done
