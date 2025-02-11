#!/bin/sh

TESTCONTAINERS_RYUK_DISABLED=true go test -v -race -buildvcs ./...

exit 0
