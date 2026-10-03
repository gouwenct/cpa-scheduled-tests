#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")"
# Native Linux and macOS builds use the same source and release packaging.
python3 build-release.py
