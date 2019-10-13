#!/usr/bin/env bash

go get github.com/gogo/protobuf/protoc-gen-gogoslick

protoc -I=. -I=$GOPATH/src -I=$GOPATH/src/github.com/gogo/protobuf/protobuf --gogoslick_out=\
Mgoogle/protobuf/struct.proto=github.com/gogo/protobuf/types:. \
./example.proto