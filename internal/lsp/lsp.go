package lsp

import (
	"encoding/json"
	"fmt"
)

var requestMethodUnmarshallers = map[string]func([]byte) (any, error){
	"initialize": unmarshallerFor[InitializeParams],
}

func unmarshallerFor[T any](data []byte) (any, error) {
	var params T
	if err := json.Unmarshal(data, &params); err != nil {
		return nil, fmt.Errorf("unmarshal %T: %w", (*T)(nil), err)
	}
	return params, nil
}

type InitializeParams struct {
	ProcessID             *Integer    `json:"processId"`
	ClientInfo            *ClientInfo `json:"clientInfo"`
	Locale                *string     `json:"locale"`
	InitializationOptions any         `json:"initializationOptions"`
	Capabilities          any         `json:"capabilities"`
	Trace                 *TraceValue `json:"trace"`
}

type ClientInfo struct {
	Name    string  `json:"name"`
	Version *string `json:"version"`
}

const (
	TraceValueOff      TraceValue = "off"
	TraceValueMessages TraceValue = "messages"
	TraceValueVerbose  TraceValue = "verbose"
)

type TraceValue string

func (t *TraceValue) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	switch TraceValue(s) {
	case TraceValueOff:
		*t = TraceValueOff
	case TraceValueMessages:
		*t = TraceValueMessages
	case TraceValueVerbose:
		*t = TraceValueVerbose
	default:
		return fmt.Errorf("unknown TraceValue: %q", s)
	}
	return nil
}

type InitializeResult struct {
	Capabilities any         `json:"capabilities"`
	ServerInfo   *ServerInfo `json:"serverInfo,omitempty"`
}

type ServerInfo struct {
	Name    string  `json:"name"`
	Version *string `json:"version"`
}
