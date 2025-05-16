#!/bin/bash

export PATH="$PATH:$(go env GOPATH)/bin"
protoc --go_out=plugins=grpc:. pb/*.proto
