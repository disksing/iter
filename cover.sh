#!/bin/bash

go test -count 10000 -coverprofile=coverage.txt -covermode=atomic
gcov2lcov -infile coverage.txt -outfile lcov.info