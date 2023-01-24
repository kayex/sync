#!/bin/sh

GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -o sync-linux-amd64 ./cmd/sync
echo 'Built sync-linux-amd'
