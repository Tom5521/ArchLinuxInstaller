#!/bin/bash

export GOARCH=arm64
go build -o builds/arm64/ArchInstaller-arm64 .
export GOARCH=amd64
go build -o builds/amd64/ArchInstaller-amd64 .
export GOARCH=386
go build -o builds/i386/ArchInstaller-i386 .
