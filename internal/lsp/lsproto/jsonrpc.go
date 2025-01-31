package lsproto

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Integer = int32

type Uinteger = uint32

type DocumentURI string // !!!

type URI string // !!!

type JSONRPCVersion struct{}

const jsonRPCVersion = `"2.0"`

func (JSONRPCVersion) MarshalJSON() ([]byte, error) {
	return []byte(jsonRPCVersion), nil
}

var ErrInvalidJSONRPCVersion = errors.New("invalid JSON-RPC version")

func (*JSONRPCVersion) UnmarshalJSON(data []byte) error {
	if string(data) != jsonRPCVersion {
		return ErrInvalidJSONRPCVersion
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

	var params any
	var retErr error

	if unmarshalParams, ok := requestMethodUnmarshallers[raw.Method]; ok {
		var err error
		params, err = unmarshalParams(raw.Params)
		if err != nil {
			retErr = fmt.Errorf("%w: %w", ErrInvalidRequest, err)
		}
	} else {
		retErr = fmt.Errorf("%w: %s", ErrMethodNotFound, raw.Method)
		var v any
		_ = json.Unmarshal(raw.Params, &v)
		params = v
	}

	r.ID = raw.ID
	r.Method = raw.Method
	r.Params = params

	return retErr
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
