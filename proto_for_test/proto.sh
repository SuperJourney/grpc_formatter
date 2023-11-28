#/bin/bash

protoc -I. --go_out=plugins=grpc:. demo.proto
protoc -I. --grpc-gateway_out=./ demo.proto