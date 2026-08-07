package ipc_test

import (
	"context"
	"net"
	"testing"

	"github.com/microsoft/typescript-go/internal/ipc"
	"github.com/microsoft/typescript-go/internal/json"
	"gotest.tools/v3/assert"
)

type noOpHandler struct{}

func (noOpHandler) HandleRequest(context.Context, string, json.Value) (any, error) {
	return nil, nil
}

func (noOpHandler) HandleNotification(context.Context, string, json.Value) error {
	return nil
}

func TestAsyncConnCallReturnsWhenPeerCloses(t *testing.T) {
	t.Parallel()
	client, server := net.Pipe()
	conn := ipc.NewAsyncConn(client, noOpHandler{})
	runDone := make(chan error, 1)
	go func() { runDone <- conn.Run(t.Context()) }()

	callDone := make(chan error, 1)
	go func() {
		_, err := conn.Call(t.Context(), "transform", nil)
		callDone <- err
	}()

	buffer := make([]byte, 1024)
	_, err := server.Read(buffer)
	assert.NilError(t, err)
	assert.NilError(t, server.Close())
	assert.NilError(t, <-runDone)
	assert.ErrorContains(t, <-callDone, "connection closed before response")
	assert.NilError(t, client.Close())
}
