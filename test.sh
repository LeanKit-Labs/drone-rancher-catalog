#!/bin/bash
cd docker
go test

cd ../github
go test

cd ../tag
go test