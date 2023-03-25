#!/bin/bash

# ----------------------------------------------------------------------------
# Integration test
# yaml_test.sh
# - Creates a talbot scaffold from YAML file
# ----------------------------------------------------------------------------
# Fail on errors
set -e

# Remove artifacts of previous builds
rm -rf testing-output
# rm ./talbot

# Rebuild binary
go build

# Scaffold a new microservice
echo "appName: testing-output
directory: .
modName: github.com/rohitkochhar/talbot-output" >> test.yaml

./talbot make -c test.yaml

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
rm -rf test.yaml
