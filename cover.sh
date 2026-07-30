#!/usr/bin/env bash

set -euo pipefail

go test -covermode=atomic -coverprofile=coverage.txt ./...
go tool cover -func=coverage.txt
