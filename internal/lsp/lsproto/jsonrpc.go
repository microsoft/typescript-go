package lsproto

import (
	"encoding/json"
	"fmt"
)

type Integer = int32

type Uinteger = uint32

type DocumentURI string // !!!

type URI string // !!!

type JSONRPCVersion struct{}

const jsonRPCVersion = "2.0"

func (JSONRPCVersion) MarshalJSON() ([]byte, error) {
	return []byte(`"` + jsonRPCVersion + `"`), nil
}

func (*JSONRPCVersion) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s != jsonRPCVersion {
		return fmt.Errorf("invalid JSON-RPC version %s", s)
	}
	return nil
}

type ID struct {
	str string
	int int32
}

func (id *ID) MarshalJSON() ([]byte, error) {
	if id.str != "" {
		return json.Marshal(id.str)
	}
	return json.Marshal(id.int)
}

func (id *ID) UnmarshalJSON(data []byte) error {
	*id = ID{}
	if len(data) > 0 && data[0] == '"' {
		return json.Unmarshal(data, &id.str)
	}
	return json.Unmarshal(data, &id.int)
}

// TODO(jakebailey): NotificationMessage? Use RequestMessage without ID?

type RequestMessage struct {
	JSONRPC JSONRPCVersion `json:"jsonrpc"`
	ID      *ID            `json:"id"`
	Method  Method         `json:"method"`
	Params  any            `json:"params"`
}

func (r *RequestMessage) UnmarshalJSON(data []byte) error {
	var raw struct {
		JSONRPC JSONRPCVersion  `json:"jsonrpc"`
		ID      *ID             `json:"id"`
		Method  Method          `json:"method"`
		Params  json.RawMessage `json:"params"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidRequest, err)
	}

	unmarshalParams, ok := requestMethodUnmarshallers[raw.Method]
	if !ok {
		return fmt.Errorf("%w: %s", ErrMethodNotFound, raw.Method)
	}

	params, err := unmarshalParams(raw.Params)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidRequest, err)
	}

	r.ID = raw.ID
	r.Method = raw.Method
	r.Params = params

	return nil
}

type ResponseMessage struct {
	JSONRPC JSONRPCVersion `json:"jsonrpc"`
	ID      *ID            `json:"id,omitempty"`
	Result  *any           `json:"result"`
	Error   *ResponseError `json:"error,omitempty"`
}

type ResponseError struct {
	Code    Integer `json:"code"`
	Message string  `json:"message"`
	Data    any     `json:"data,omitempty"`
}
