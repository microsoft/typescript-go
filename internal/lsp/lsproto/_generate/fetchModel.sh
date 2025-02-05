#!/usr/bin/env bash

cd "$(dirname "$0")"

hash=dadd73f7fc283b4d0adb602adadcf4be16ef3a7b

curl -sL "https://raw.githubusercontent.com/microsoft/vscode-languageserver-node/$hash/protocol/metaModel.json" > metaModel.json
curl -sL "https://raw.githubusercontent.com/microsoft/vscode-languageserver-node/$hash/protocol/metaModel.schema.json" > metaModel.schema.json

json -I -f metaModel.schema.json -e 'Object.assign(this, this.definitions.MetaModel); delete this.definitions.MetaModel'

npx json-schema-to-typescript metaModel.schema.json > metaModelSchema.mts
