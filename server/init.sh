#!/bin/bash

cd gRPC
protoc --go_out=../backpackTaskGRPC/ --go_opt=paths=source_relative --go-grpc_out=../backpackTaskGRPC/ --go-grpc_opt=paths=source_relative grpc.proto
cd ..
go build
exec ./server
