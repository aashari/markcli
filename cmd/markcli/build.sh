#!/bin/bash
set -e

echo "Starting the build process"

# Create a dist directory to store builds
mkdir -p dist

# Build executables for different platforms
echo "Building for Linux"
GOOS=linux GOARCH=amd64 go build -o dist/markcli-linux ./cmd/markcli

echo "Building for macOS"
GOOS=darwin GOARCH=amd64 go build -o dist/markcli-mac ./cmd/markcli

echo "Building for Windows"
GOOS=windows GOARCH=amd64 go build -o dist/markcli.exe ./cmd/markcli

echo "Build process completed." 