#!/usr/bin/env bash
flags="-s -w"
export CGO_ENABLED=0
go build -ldflags="$flags"
GOOS=windows GOARCH=amd64 go build -ldflags="$flags"
