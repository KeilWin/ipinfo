#!/bin/bash

PROJECT_NAME="ipinfo"

cp configs/$PROJECT_NAME.env bin/configs/$PROJECT_NAME.env

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -o bin/$PROJECT_NAME \
    -trimpath \
    -ldflags="-s -w" \
    ./cmd/$PROJECT_NAME
