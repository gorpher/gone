#! /bin/bash
basepath=$(pwd)
echo "now at ${basepath}"
golangci-lint run -c ./.golangci.yml
go mod tidy
