package lsp

import (
	"fmt"

	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/diagnostics"
	"github.com/microsoft/typescript-go/internal/jsonrpc"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
)

type progressEvent struct {
	message *diagnostics.Message
	args    []any
	finish  bool
}

// projectLoadingProgress manages LSP WorkDoneProgress indicators for project
// loading. A single persistent goroutine processes start/finish events,
// maintains the set of loading projects, and sends progress messages in order.
//
// Callers never block: events are sent to a buffered channel with a select
// on the caller's context.
type projectLoadingProgress struct {
	server *Server
	ch     chan progressEvent
}

func newProjectLoadingProgress(server *Server) *projectLoadingProgress {
	p := &projectLoadingProgress{
		server: server,
		ch:     make(chan progressEvent, 64),
	}
	go p.run()
	return p
}

func (p *projectLoadingProgress) start(message *diagnostics.Message, args ...any) {
	select {
	case p.ch <- progressEvent{message: message, args: args}:
		// Sent successfully.
	case <-p.server.backgroundCtx.Done():
		// Server shutting down; drop the event.
	}
}

func (p *projectLoadingProgress) finish(message *diagnostics.Message, args ...any) {
	select {
	case p.ch <- progressEvent{message: message, args: args, finish: true}:
		// Sent successfully.
	case <-p.server.backgroundCtx.Done():
		// Server shutting down; drop the event.
	}
}

// run is the persistent goroutine that processes all progress events.
// It owns all mutable state: no external synchronization needed.
func (p *projectLoadingProgress) run() {
	var (
		loading collections.OrderedSet[string]
		token   string // current token; empty if no progress active
		tokenID int
		begun   bool // whether "begin" has been sent for the current token
		title   = diagnostics.Loading.Localize(p.server.locale)
	)

	localize := func(ev progressEvent) string {
		return ev.message.Localize(p.server.locale, ev.args...)
	}

	for {
		select {
		case ev := <-p.ch:
			text := localize(ev)
			if !ev.finish {
				loading.Add(text)
				if token == "" {
					tokenID++
					token = fmt.Sprintf("tsgo-loading-%d", tokenID)
					begun = false
					p.sendProgressCreate(token)
				}
				if !begun {
					begun = true
					p.sendProgress(token, &lsproto.WorkDoneProgressBegin{
						Title:   title,
						Message: &text,
					})
				} else {
					p.sendProgress(token, &lsproto.WorkDoneProgressReport{
						Message: &text,
					})
				}
			} else {
				loading.Delete(text)
				if token == "" {
					continue
				}
				if loading.Size() == 0 {
					p.sendProgress(token, &lsproto.WorkDoneProgressEnd{})
					token = ""
				} else {
					first := firstValue(loading.Values())
					p.sendProgress(token, &lsproto.WorkDoneProgressReport{
						Message: &first,
					})
				}
			}

		case <-p.server.backgroundCtx.Done():
			return
		}
	}
}

// sendProgress sends a $/progress notification with a snapshot of the token
// string, so deferred serialization in the write loop won't see a mutated value.
func (p *projectLoadingProgress) sendProgress(token string, value any) {
	_ = sendNotification(p.server, lsproto.ProgressInfo, &lsproto.ProgressParams{
		Token: lsproto.IntegerOrString{String: &token},
		Value: value,
	})
}

// sendProgressCreate sends a window/workDoneProgress/create request without
// waiting for the response. The client processes messages in order, so
// subsequent $/progress notifications will arrive after the create.
func (p *projectLoadingProgress) sendProgressCreate(token string) {
	id := jsonrpc.NewIDString(fmt.Sprintf("ts%d", p.server.clientSeq.Add(1)))
	req := lsproto.WindowWorkDoneProgressCreateInfo.NewRequestMessage(id, &lsproto.WorkDoneProgressCreateParams{
		Token: lsproto.IntegerOrString{String: &token},
	})
	_ = p.server.send(req.Message())
}

func firstValue[T any](seq func(yield func(T) bool)) T {
	for v := range seq {
		return v
	}
	panic("empty sequence")
}
