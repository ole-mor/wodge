#!/bin/bash

# Build the wodge binary
echo "Building wodge..."
go build -o wodge cmd/wodge/main.go

if [ $? -ne 0 ]; then
    echo "Build failed."
    exit 1
fi

# Move to /usr/local/bin (requires sudo)
echo "Installing wodge to /usr/local/bin..."
if [ -w /usr/local/bin ]; then
    mv wodge /usr/local/bin/
else
    sudo mv wodge /usr/local/bin/
fi

if [ $? -eq 0 ]; then
    echo "Successfully installed wodge!"
    echo "Run 'wodge help' to get started."
else
    echo "Installation failed."
    exit 1
fi
