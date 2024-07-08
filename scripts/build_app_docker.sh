#!/bin/sh

DOCKER_TAG=prulloac/fineasy:latest
DOCKER_ARGS=--no-cache

docker build -t ${DOCKER_TAG} -f deployment/Dockerfile ${DOCKER_ARGS} .

exit 0
