package api

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"
	"sync/atomic"

	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
	"github.com/microsoft/typescript-go/internal/jsonrpc"
)

var (
	ErrConnClosed     = errors.New("api: connection closed")
	ErrRequestTimeout = errors.New("api: request timeout")
)

// Handler processes incoming API requests and notifications.
type Handler interface {
	// HandleRequest handles an incoming request and returns a result or error.
	HandleRequest(ctx context.Context, method string, params jsontext.Value) (any, error)
	// HandleNotification handles an incoming notification.
	HandleNotification(ctx context.Context, method string, params jsontext.Value) error
}

// Conn manages bidirectional JSON-RPC communication over a connection.
type Conn struct {
	rwc      io.ReadWriteCloser
	protocol Protocol
	handler  Handler

	// For server→client requests
	seq        atomic.Int64
	pending    map[jsonrpc.ID]chan *Message
	pendingMu  sync.Mutex
	writeMu    sync.Mutex
	closed     atomic.Bool
	closedChan chan struct{}
}

// NewConn creates a new connection with the given transport and handler.
func NewConn(rwc io.ReadWriteCloser, handler Handler) *Conn {
	return &Conn{
		rwc:        rwc,
		protocol:   NewJSONRPCProtocol(rwc),
		handler:    handler,
		pending:    make(map[jsonrpc.ID]chan *Message),
		closedChan: make(chan struct{}),
	}
}

// Run starts processing messages on the connection.
// It blocks until the connection is closed or an error occurs.
func (c *Conn) Run(ctx context.Context) error {
	defer c.Close()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-c.closedChan:
			return ErrConnClosed
		default:
			// Non-blocking check - continue to read messages
		}

		msg, err := c.protocol.ReadMessage()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}

		if msg.IsResponse() {
			c.handleResponse(msg)
		} else if msg.IsRequest() {
			go c.handleRequest(ctx, msg)
		} else if msg.IsNotification() {
			go c.handleNotification(ctx, msg)
		}
	}
}

// handleResponse matches a response to a pending request.
func (c *Conn) handleResponse(msg *Message) {
	c.pendingMu.Lock()
	ch, ok := c.pending[*msg.ID]
	if ok {
		delete(c.pending, *msg.ID)
	}
	c.pendingMu.Unlock()

	if ok {
		ch <- msg
		close(ch)
	}
}

// handleRequest processes an incoming request.
func (c *Conn) handleRequest(ctx context.Context, msg *Message) {
	result, err := c.handler.HandleRequest(ctx, msg.Method, msg.Params)

	c.writeMu.Lock()
	defer c.writeMu.Unlock()

	if err != nil {
		_ = c.protocol.WriteError(msg.ID, &jsonrpc.ResponseError{
			Code:    jsonrpc.CodeInternalError,
			Message: err.Error(),
		})
		return
	}

	_ = c.protocol.WriteResponse(msg.ID, result)
}

// handleNotification processes an incoming notification.
func (c *Conn) handleNotification(ctx context.Context, msg *Message) {
	_ = c.handler.HandleNotification(ctx, msg.Method, msg.Params)
}

// Call sends a request to the client and waits for a response.
func (c *Conn) Call(ctx context.Context, method string, params any) (jsontext.Value, error) {
	if c.closed.Load() {
		return nil, ErrConnClosed
	}

	id := jsonrpc.NewIDString(fmt.Sprintf("s%d", c.seq.Add(1)))
	ch := make(chan *Message, 1)

	c.pendingMu.Lock()
	c.pending[*id] = ch
	c.pendingMu.Unlock()

	c.writeMu.Lock()
	err := c.protocol.WriteRequest(id, method, params)
	c.writeMu.Unlock()

	if err != nil {
		c.pendingMu.Lock()
		delete(c.pending, *id)
		c.pendingMu.Unlock()
		return nil, err
	}

	select {
	case <-ctx.Done():
		c.pendingMu.Lock()
		delete(c.pending, *id)
		c.pendingMu.Unlock()
		return nil, ctx.Err()
	case <-c.closedChan:
		return nil, ErrConnClosed
	case resp := <-ch:
		if resp.Error != nil {
			return nil, fmt.Errorf("api: remote error [%d]: %s", resp.Error.Code, resp.Error.Message)
		}
		return resp.Result, nil
	}
}

// Notify sends a notification to the client (no response expected).
func (c *Conn) Notify(ctx context.Context, method string, params any) error {
	if c.closed.Load() {
		return ErrConnClosed
	}

	c.writeMu.Lock()
	defer c.writeMu.Unlock()
	return c.protocol.WriteNotification(method, params)
}

// Close closes the connection.
func (c *Conn) Close() error {
	if c.closed.CompareAndSwap(false, true) {
		close(c.closedChan)

		// Cancel all pending requests
		c.pendingMu.Lock()
		for id, ch := range c.pending {
			delete(c.pending, id)
			close(ch)
		}
		c.pendingMu.Unlock()

		return c.rwc.Close()
	}
	return nil
}

// UnmarshalParams is a helper to unmarshal params into a typed struct.
func UnmarshalParams[T any](params jsontext.Value) (*T, error) {
	if len(params) == 0 {
		return nil, nil
	}
	var v T
	if err := json.Unmarshal(params, &v); err != nil {
		return nil, err
	}
	return &v, nil
}
