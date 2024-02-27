#!/bin/bash

protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=./ --go-grpc_opt=paths=source_relative ./blobcache.proto
protoc -I ./ --python_betterproto_out=..  ./blobcache.proto
