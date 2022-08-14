#!/bin/bash

mkdir -p /context/user_api
cp -r /src/user_api /context/

pushd /context/user_api
go mod tidy
go build
go install
user_api
