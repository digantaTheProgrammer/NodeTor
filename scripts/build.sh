#!/usr/bin/env bash
set -exuo pipefail

cd "$( dirname "${BASH_SOURCE[0]}" )/.."
source .envrc

GOOS=linux go build -mod=vendor -ldflags="-s -w" -o bins/supply ./src/nodejs/supply/cli
GOOS=linux go build -mod=vendor -ldflags="-s -w" -o bins/finalize ./src/nodejs/finalize/cli
