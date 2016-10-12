#!/bin/bash
#compiles project for use in an Alpine docker image (i.e. docker-dind)

export CGO_ENABLED=0 
export GOOS=linux 
export GOARCH=amd64 

echo "compiling project"
go build -a -tags netgo