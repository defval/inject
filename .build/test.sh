#!/usr/bin/env bash

set -e

if [[ -f coverage.txt ]]; then
    rm coverage.txt
fi

go test -coverprofile=profile.out -covermode=atomic ./...

if [[ -f profile.out ]]; then
    cat profile.out >> coverage.txt
    rm profile.out
fi