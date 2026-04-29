#!/bin/bash

echo "Building for Linux (amd64)..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o coupang_spider -ldflags "-s -w" coupang_spider.go

if [ $? -eq 0 ]; then
    echo "Build successful. Binary 'coupang_spider' created."
else
    echo "Build failed."
    exit 1
fi
