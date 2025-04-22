#!/bin/bash

set -exo pipefail
cd "$(dirname "$0")"

[[ -f ./tsgo ]] || go build ./cmd/tsgo
[[ -f ./tsgo.wasm ]] || GOOS=wasip1 GOARCH=wasm go build -o tsgo.wasm ./cmd/tsgo
[[ -f ./tsgo-tinygo.wasm ]] || GOOS=wasip1 GOARCH=wasm tinygo build -o tsgo-tinygo.wasm ./cmd/tsgo

hyperfine -w 1 \
    -n 'native' \
    './tsgo -p $PWD/_submodules/TypeScript/src/compiler --singleThreaded --listFilesOnly' \
    -n 'go wasm' \
    'wasmtime run --dir=/ --env PWD="$PWD" --env PATH="$PATH" -W max-wasm-stack=1048576 ./tsgo.wasm -p $PWD/_submodules/TypeScript/src/compiler --singleThreaded --listFilesOnly' \
    -n 'tinygo wasm' \
    'wasmtime run --dir=/ --env PWD="$PWD" --env PATH="$PATH" -W max-wasm-stack=1048576 ./tsgo-tinygo.wasm -p $PWD/_submodules/TypeScript/src/compiler --singleThreaded --listFilesOnly'
