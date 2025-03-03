package api

import (
	"encoding/json"
	"fmt"
)

type Method string

const (
	MethodParseConfigFile     Method = "parseConfigFile"
	MethodLoadProject         Method = "loadProject"
	MethodGetSymbolAtPosition Method = "getSymbolAtPosition"
)

var unmarshalers = map[Method]func([]byte) (any, error){
	MethodParseConfigFile:     unmarshallerFor[ParseConfigFileParams],
	MethodLoadProject:         unmarshallerFor[LoadProjectParams],
	MethodGetSymbolAtPosition: unmarshallerFor[GetSymbolAtPositionParams],
}

type ParseConfigFileParams struct {
	FileName string `json:"fileName"`
}

type LoadProjectParams struct {
	ConfigFileName string `json:"configFileName"`
}

type GetSymbolAtPositionParams struct {
	Project  string `json:"project"`
	FileName string `json:"fileName"`
	Position uint32 `json:"position"`
}

func unmarshalPayload(method string, payload json.RawMessage) (any, error) {
	unmarshaler, ok := unmarshalers[Method(method)]
	if !ok {
		return nil, fmt.Errorf("unknown API method %q", method)
	}
	return unmarshaler(payload)
}

func unmarshallerFor[T any](data []byte) (any, error) {
	var v T
	if err := json.Unmarshal(data, &v); err != nil {
		return nil, fmt.Errorf("failed to unmarshal %T: %w", (*T)(nil), err)
	}
	return &v, nil
}
