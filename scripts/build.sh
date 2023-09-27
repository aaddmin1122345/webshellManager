#!/bin/bash

if [ -f "go.mod" ]; then
  go build -ldflags "-s -w" -o main ./cmd/webshellManager/main.go
else
  go build -modfile ../go.mod -ldflags "-s -w" -o main ../cmd/webshellManager/main.go
fi
