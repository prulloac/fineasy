#!/bin/sh

go fmt ./...
go mod tidy -v

exit 0
