#!/bin/sh

go test -v -race -buildvcs -coverprofile=/tmp/coverage.out ./...
go tool cover -html=/tmp/coverage.out

exit 0
