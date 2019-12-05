#!/bin/bash

go test -count 10000 -coverprofile=coverage.txt -covermode=atomic
sed -i '/reflection.go/d' coverage.txt
gcov2lcov -infile coverage.txt -outfile lcov.info
