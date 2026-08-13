package lsp

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/microsoft/typescript-go/internal/jsonrpc"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
)

type eofReader struct{}

func (eofReader) Read() (*lsproto.Message, error) { return nil, io.EOF }

// A response that exceeds the JSON encoder's nesting limit must fail only its
// request. The write loop must remain available to deliver subsequent responses.
func TestWriteLoopRecoversFromUnserializableResponse(t *testing.T) {
	t.Parallel()

	pr, pw := io.Pipe()
	server := NewServer(&ServerOptions{
		In:  eofReader{},
		Out: ToWriter(pw),
		Err: io.Discard,
		Cwd: "/test",
	})

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	server.backgroundCtx = ctx

	writeLoopErr := make(chan error, 1)
	go func() { writeLoopErr <- server.writeLoop(ctx) }()

	// A selection range whose parent chain is far deeper than the JSON encoder's nesting limit.
	var deep *lsproto.SelectionRange
	for range 20000 {
		deep = &lsproto.SelectionRange{Parent: deep}
	}
	badResult := []*lsproto.SelectionRange{deep}
	badID := jsonrpc.NewIDString("bad")
	if err := server.send((&lsproto.ResponseMessage{ID: badID, Result: &badResult}).Message()); err != nil {
		t.Fatalf("failed to enqueue bad response: %v", err)
	}

	// A subsequent well-formed response must still be delivered.
	goodID := jsonrpc.NewIDString("good")
	if err := server.send((&lsproto.ResponseMessage{ID: goodID, Result: &lsproto.SelectionRangesOrNull{}}).Message()); err != nil {
		t.Fatalf("failed to enqueue good response: %v", err)
	}

	reader := lsproto.NewBaseReader(pr)
	sawError := false
	sawGood := false
	for range 2 {
		msg := readMessageWithTimeout(t, reader)
		resp := msg.AsResponse()
		switch {
		case resp.ID != nil && *resp.ID == *badID:
			if resp.Error == nil {
				t.Errorf("expected an error response for the unserializable request, got a result")
			} else if resp.Error.Code != int32(lsproto.ErrorCodeInternalError) {
				t.Errorf("error response code = %d, want %d", resp.Error.Code, lsproto.ErrorCodeInternalError)
			}
			sawError = true
		case resp.ID != nil && *resp.ID == *goodID:
			if resp.Error != nil {
				t.Errorf("expected a successful response for the good request, got error: %v", resp.Error)
			}
			sawGood = true
		default:
			t.Errorf("unexpected response id: %v", resp.ID)
		}
	}

	if !sawError {
		t.Errorf("did not receive an error response for the unmarshalable request")
	}
	if !sawGood {
		t.Errorf("did not receive the subsequent well-formed response (write loop likely died)")
	}

	// The write loop must still be running.
	select {
	case err := <-writeLoopErr:
		t.Fatalf("write loop exited unexpectedly: %v", err)
	default:
		return
	}
}

func readMessageWithTimeout(t *testing.T, reader *lsproto.BaseReader) *lsproto.Message {
	t.Helper()
	type result struct {
		msg *lsproto.Message
		err error
	}
	ch := make(chan result, 1)
	go func() {
		data, err := reader.Read()
		if err != nil {
			ch <- result{err: err}
			return
		}
		msg := &lsproto.Message{}
		ch <- result{msg: msg, err: msg.UnmarshalJSON(data)}
	}()
	select {
	case r := <-ch:
		if r.err != nil {
			t.Fatalf("failed to read message: %v", r.err)
		}
		return r.msg
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for a message (write loop may have died)")
		return nil
	}
}
