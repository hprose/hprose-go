#!/usr/bin/env bash

for d in $(go list ./... | grep -v vendor); do
    go test $d
done