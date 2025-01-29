package lsp

import (
	"encoding/json"
)

type Integer int32

type Uinteger int32

type DocumentURI string

type URI string

type ID struct {
	str string
	int int32
}

func (id ID) MarshalJSON() ([]byte, error) {
	if id.str != "" {
		return json.Marshal(id.str)
	}
	return json.Marshal(id.int)
}

func (id *ID) UnmarshalJSON(data []byte) error {
	*id = ID{}
	if data[0] == '"' {
		return json.Unmarshal(data, &id.str)
	}
	return json.Unmarshal(data, &id.int)
}

type RequestParams interface {
	requestParams()
}

type RequestMessage struct {
	JSONRPC string        `json:"jsonrpc"`
	ID      ID            `json:"id"`
	Method  string        `json:"method"`
	Params  RequestParams `json:"params"`
}

func (r *RequestMessage) UnmarshalJSON(data []byte) error {
	var raw struct {
		JSONRPC string          `json:"jsonrpc"`
		ID      ID              `json:"id"`
		Method  string          `json:"method"`
		Params  json.RawMessage `json:"params"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	unmarshalParams, ok := requestMethodUnmarshallers[raw.Method]
	if !ok {
		// TODO: use a real error
		return fmt.Errorf("unknown method %s", raw.Method)
	}

	
}
