#!/bin/bash



for file in $(find . -name '*.go' | grep -v /vendor/); do
    if grep -q '^type.*interface {$' ${file}; then
        relative_path=${file#./internal/}
        dir_path=$(dirname ${relative_path})
        if [[ ${dir_path} == "." ]]; then
            dest=$(basename ${file} .go)_mock.go
        else
            dest="${dir_path}/$(basename ${file} .go)_mock.go"
        fi
        mockgen -source=${file} -destination=test/mock/${dest}
    fi
done
