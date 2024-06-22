#!/bin/sh

MAIN_PACKAGE_PATH=./cmd/app
BINARY_NAME=fineasy
OUTPUT_DIR=./bin

# Include additional build steps, like TypeScript, SCSS or Tailwind compilation here...
go build -o=${OUTPUT_DIR}/${BINARY_NAME} ${MAIN_PACKAGE_PATH}

exit 0
