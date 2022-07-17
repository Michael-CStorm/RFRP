#!/bin/bash

pushd /src/user_api
go build
go install
user_api
