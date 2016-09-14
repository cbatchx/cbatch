#!/bin/bash
set -x

version=0.0.6

export CGO_ENABLED=1
export GOOS=linux
export GOARCH=amd64

cd ../ && go build

mkdir -p release

tar -cvf release/cbatch-$version-amd64linux.tar cbatch config/config.toml config/bootstrap.tmpl.sh 