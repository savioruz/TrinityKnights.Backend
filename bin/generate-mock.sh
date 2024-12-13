#!/bin/bash

for file in $(find . -name '*.go' | grep -v /vendor/); do
    if grep -q '^type.*interface {$' ${file}; then
        dest=$(basename ${file} .go)_mock.go
        echo "Creating mock for ${file} -> ${dest}"
        cat <<EOF > test/mock/${dest}
package mock

import (
    "github.com/stretchr/testify/mock"
)

// Auto-generated mock for ${file}
EOF
    fi
done
