#!/bin/sh
go get -u github.com/tidwall/gjson
mkdir -p artifacts
go build -ldflags "-s -w" -o artifacts/psguard