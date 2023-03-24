#!/bin/bash

# ----------------------------------------------------------------------------
# Integration test
# ----------------------------------------------------------------------------
# Fail on errors
set -e

# Remove artifacts of previous builds
rm -rf testing-output

# Scaffold a new microservice
go run main.go make -n testing-output -d . -m github.com/rohitkochhar/talbot-output

cd ./testing-output

# Run automatically generated unit tests
cd ./cmd/api
go test -v ./
cd ../../

# Check if we can build docker image
docker build -t talbot:testing-output .

# Cleanup
docker rmi talbot:testing-output
cd ../
rm -rf testing-output
