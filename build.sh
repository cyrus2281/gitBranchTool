#!/bin/bash

# Version file path 
FILE_PATH="./internal/version.go"

# Extract the version string using grep and sed
VERSION=$(grep 'const VERSION =' "$FILE_PATH" | sed -E 's/.*"([^"]+)".*/\1/')

# Print the version string
echo "Building for version $VERSION"

# Create bin directory
rm -rf bin
mkdir -p bin

# Build for all platforms
echo "Building for linux..."
GOOS=linux go build -o bin/g-linux-v$VERSION

echo "Building for macos..."
GOOS=darwin go build -o bin/g-macos-v$VERSION

echo "Building for windows..."
GOOS=windows go build -o bin/g-win-v$VERSION.exe

echo "Done building for all platforms."
echo -e "\tBuild directory: bin/"

