package lsproto

import (
	"encoding/json"
	"errors"
	"fmt"
)

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

type Message struct {
	JSONRPC JSONRPCVersion `json:"jsonrpc"`
}

type NotificationMessage struct {
	Message
	Method Method `json:"method"`
	Params any    `json:"params"`
}

type RequestMessage struct {
	Message
	ID     *ID    `json:"id"`
	Method Method `json:"method"`
	Params any    `json:"params"`
}

type RequestOrNotificationMessage struct {
	NotificationMessage *NotificationMessage
	RequestMessage      *RequestMessage
}

func (r *RequestOrNotificationMessage) UnmarshalJSON(data []byte) error {
	var raw struct {
		JSONRPC JSONRPCVersion  `json:"jsonrpc"`
		ID      *ID             `json:"id,omitzero"`
		Method  Method          `json:"method"`
		Params  json.RawMessage `json:"params"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidRequest, err)
	}

	params, err := unmarshalParams(raw.Method, raw.Params)
	if err != nil {
		return err
	}

	if raw.ID != nil {
		r.RequestMessage = &RequestMessage{
			ID:     raw.ID,
			Method: raw.Method,
			Params: params,
		}
	} else {
		r.NotificationMessage = &NotificationMessage{
			Method: raw.Method,
			Params: params,
		}
	}

	return nil
}

func unmarshalParams(rawMethod Method, rawParams []byte) (any, error) {
	if rawMethod == MethodShutdown || rawMethod == MethodExit {
		// These methods have no params.
		return nil, nil
	}

	var params any
	var err error

	if unmarshaller, ok := unmarshallers[rawMethod]; ok {
		params, err = unmarshaller(rawParams)
	} else {
		// Fall back to default; it's probably an unknown message and we will probably not handle it.
		err = json.Unmarshal(rawParams, &params)
	}

	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidRequest, err)
	}

	return params, nil
}

type ResponseMessage struct {
	Message
	ID     *ID            `json:"id,omitzero"`
	Result any            `json:"result"`
	Error  *ResponseError `json:"error,omitzero"`
}

type ResponseError struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitzero"`
}
