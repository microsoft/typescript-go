package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/microsoft/typescript-go/_tools/genlsp/metamodel"
)

const metaModelURL = `https://raw.githubusercontent.com/microsoft/vscode-languageserver-node/dadd73f7fc283b4d0adb602adadcf4be16ef3a7b/protocol/metaModel.json`

func main() {
	model := getModel()

	for _, t := range model.Requests {
		fmt.Println(t.Method)
	}
}

var supportedMethods = map[string]bool{
	"initialize":  true,
	"shutdown":    true,
	"exit":        true,
	"initialized": true,
}

func getModel() *metamodel.MetaModel {
	resp, err := http.Get(metaModelURL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var model metamodel.MetaModel
	if err := json.NewDecoder(resp.Body).Decode(&model); err != nil {
		panic(err)
	}

	return &model
}
