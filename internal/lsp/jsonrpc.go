package lsp

import (
	"encoding/json"
	"fmt"
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

type RequestMessage struct {
	JSONRPC string `json:"jsonrpc"`
	ID      ID     `json:"id"`
	Method  string `json:"method"`
	Params  any    `json:"params"`
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
