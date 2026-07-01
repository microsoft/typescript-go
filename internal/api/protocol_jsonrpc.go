package api

import (
	"io"

	"github.com/microsoft/typescript-go/internal/json"
	"github.com/microsoft/typescript-go/internal/jsonrpc"
)

// JSONRPCProtocol implements the Protocol interface using JSON-RPC 2.0
// with the LSP base protocol framing (Content-Length headers).
type JSONRPCProtocol struct {
	reader *jsonrpc.Reader
	writer *jsonrpc.Writer
}

var _ Protocol = (*JSONRPCProtocol)(nil)

// NewJSONRPCProtocol creates a new JSON-RPC protocol handler.
func NewJSONRPCProtocol(rw io.ReadWriter) *JSONRPCProtocol {
	return &JSONRPCProtocol{
		reader: jsonrpc.NewReader(rw),
		writer: jsonrpc.NewWriter(rw),
	}
}

// ReadMessage implements Protocol.
func (p *JSONRPCProtocol) ReadMessage() (*Message, error) {
	data, err := p.reader.Read()
	if err != nil {
		return nil, err
	}

	var msg Message
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, err
	}

	return &msg, nil
}

// WriteRequest implements Protocol.
func (p *JSONRPCProtocol) WriteRequest(id *jsonrpc.ID, method string, params any) error {
	msg := jsonrpc.RequestMessage{
		ID:     id,
		Method: method,
		Params: params,
	}
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return p.writer.Write(data)
}

// WriteNotification implements Protocol.
func (p *JSONRPCProtocol) WriteNotification(method string, params any) error {
	msg := jsonrpc.RequestMessage{
		Method: method,
		Params: params,
	}
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return p.writer.Write(data)
}

// jsonrpcTimedResult is the envelope used to carry server processing-time
// metadata over JSON-RPC. Because the vscode-jsonrpc client library surfaces
// only the response's `result` member to request callers, timing is embedded
// within `result` rather than added as a sibling member. The client unwraps
// this envelope when it has requested timing collection.
type jsonrpcTimedResult struct {
	// Result is the underlying handler result.
	Result any `json:"result"`
	// ServerTimeMicros is the server-side processing time in microseconds.
	ServerTimeMicros uint32 `json:"serverTimeMicros"`
}

// WriteResponse implements Protocol.
func (p *JSONRPCProtocol) WriteResponse(id *jsonrpc.ID, result any) error {
	// When timing is enabled, wrap the result in an envelope carrying the
	// server's processing time. The client unwraps it to recover the result.
	if timed, ok := result.(serverTimedResult); ok {
		inner := timed.result
		if inner == nil {
			inner = json.Value("null")
		}
		result = jsonrpcTimedResult{
			Result:           inner,
			ServerTimeMicros: timed.processingTimeMicros,
		}
	}
	if result == nil {
		result = json.Value("null")
	}
	msg := jsonrpc.ResponseMessage{
		ID:     id,
		Result: result,
	}
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return p.writer.Write(data)
}

// WriteError implements Protocol.
func (p *JSONRPCProtocol) WriteError(id *jsonrpc.ID, respErr *jsonrpc.ResponseError) error {
	msg := jsonrpc.ResponseMessage{
		ID:    id,
		Error: respErr,
	}
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return p.writer.Write(data)
}
