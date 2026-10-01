#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")"
mkdir -p dist
CGO_ENABLED=1 go build -trimpath -buildmode=c-shared -o dist/cpa-scheduled-tests.so .
rm -f dist/cpa-scheduled-tests.h
printf 'Built: %s\n' "$(pwd)/dist/cpa-scheduled-tests.so"
