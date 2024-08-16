#!/usr/bin/env bash
version=$(git describe --abbrev=0 --tags)
flags="-s -w -X main.Version=$version"
export CGO_ENABLED=0
go build -ldflags="$flags"
GOOS=windows GOARCH=amd64 go build -ldflags="$flags"
