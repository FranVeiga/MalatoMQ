#!/bin/bash

if [ $# -ne 1 ]; then
    echo "Incorrect number of arguments";
    exit 1
fi

protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative $1
