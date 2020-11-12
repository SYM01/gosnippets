#!/bin/bash

GOOS=darwin GOARCH=amd64 go build -o output/curl_darwin -ldflags "-s -w" .

GOOS=linux GOARCH=amd64 go build -o output/curl_linux -ldflags "-s -w" .
