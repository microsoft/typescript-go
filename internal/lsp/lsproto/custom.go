package lsproto

import (
	"encoding/json"
	"fmt"
)

const (
	MethodAPI Method = "custom/api"
)

type APIRequestParams struct {
	Method APIMethod `json:"method"`
	Params any       `json:"params"`
}

type APIMethod string

const (
	APIMethodParseConfigFile     APIMethod = "parseConfigFile"
	APIMethodLoadProject         APIMethod = "loadProject"
	APIMethodGetSymbolAtPosition APIMethod = "getSymbolAtPosition"
)

var apiUnmarshalers = map[APIMethod]func([]byte) (any, error){
	APIMethodParseConfigFile:     unmarshallerFor[APIParseConfigFileParams],
	APIMethodLoadProject:         unmarshallerFor[APILoadProjectParams],
	APIMethodGetSymbolAtPosition: unmarshallerFor[APIGetSymbolAtPositionParams],
}

type APIParseConfigFileParams struct {
	FileName string `json:"fileName"`
}

type APILoadProjectParams struct {
	ConfigFileName string `json:"configFileName"`
}

type APIGetSymbolAtPositionParams struct {
	Project  string `json:"project"`
	FileName string `json:"fileName"`
	Position uint32 `json:"position"`
}

func (p *APIRequestParams) UnmarshalJSON(data []byte) error {
	*p = APIRequestParams{}
	var raw struct {
		Method APIMethod       `json:"method"`
		Params json.RawMessage `json:"params"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	p.Method = raw.Method

	var params any
	var err error

	if unmarshaler, ok := apiUnmarshalers[raw.Method]; ok {
		params, err = unmarshaler(raw.Params)
	} else {
		return fmt.Errorf("unknown API method %q", raw.Method)
	}
	p.Params = params

	if err != nil {
		return fmt.Errorf("failed to unmarshal API method %q: %w", raw.Method, err)
	}
	return nil
}
