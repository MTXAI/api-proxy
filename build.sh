#!/bin/bash

set -x

rm -rf output
mkdir output

GOOS=linux GOARCH=amd64 go build -v -o ./output/ ./cmd/...
