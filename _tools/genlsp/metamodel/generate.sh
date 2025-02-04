#!/usr/bin/env bash

cd "$(dirname "$0")"

curl -sL https://raw.githubusercontent.com/microsoft/vscode-languageserver-node/dadd73f7fc283b4d0adb602adadcf4be16ef3a7b/protocol/metaModel.schema.json | go run github.com/atombender/go-jsonschema@latest -p metamodel --tags json - > metamodel_generated.go
