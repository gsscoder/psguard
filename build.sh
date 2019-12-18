#!/bin/sh
mkdir -p artifacts
go build -ldflags "-s -w" -o artifacts/psguard