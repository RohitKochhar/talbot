#!/bin/bash

# Parse flags
while getopts n:d:m: flag
do case "${flag}" in 
    n) name=${OPTARG};;
    d) dir=${OPTARG};;
    m) modname=${OPTARG};;
    esac
done

# Check if we need to set $dir to default
if [ -z "$dir" ]
then
    dir="./"
fi

# Check if a name was provided, fail otherwise
if [ -z "$name" ]
then
    echo "Usage: skele-server -n APP_NAME -m MOD_NAME -d DIRECTORY (default: ./)"
    exit 0
fi

# Check if a mod name was given, otherwise set to name
if [ -z "$modname" ]
then
    modname=$name
fi

echo "Creating application $name (modname: $modname) in $dir"

cd $dir

echo "Creating subdirectory: $name"

mkdir $name

cd $name

echo "Initializing go module: $modname"

go mod init $modname

echo "Creating base files"

mkdir -p bin cmd/api internal migrations remote
touch Makefile
touch cmd/api/main.go

