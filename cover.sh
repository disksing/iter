#!/bin/bash

go test -coverprofile=coverage.txt -covermode=atomic ./algo -fuzz-time=5m
gcov2lcov -infile coverage.txt -outfile lcov.info
