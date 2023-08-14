#!/bin/bash

export GOARCH=arm64
go build -o builds/arm64/ArchInstaller-arm64 .
export GOARCH=amd64
go build -o builds/amd64/ArchInstaller-amd64 .
