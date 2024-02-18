#!/bin/bash

protoc --go_out=./internal/pb --go-grpc_out=./internal/pb ./internal/proto/*.proto 

