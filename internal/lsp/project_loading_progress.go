package lsp

import (
	"context"
	"fmt"

	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/jsonrpc"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/tspath"
)

type progressEvent struct {
	name   string // project name
	finish bool   // true if the project finished loading
}

// projectLoadingProgress manages LSP WorkDoneProgress indicators for project
// loading. A single persistent goroutine processes start/finish events,
// maintains the set of loading projects, and sends progress messages in order.
//
// Callers never block: events are sent to a buffered channel with a select
// on the server's background context.
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

func (p *projectLoadingProgress) start(ctx context.Context, projectName string) {
	select {
	case p.ch <- progressEvent{name: projectName}:
	case <-ctx.Done():
	}
}

func (p *projectLoadingProgress) finish(ctx context.Context, projectName string) {
	select {
	case p.ch <- progressEvent{name: projectName, finish: true}:
	case <-ctx.Done():
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
	)

	for {
		select {
		case ev := <-p.ch:
			if !ev.finish {
				loading.Add(ev.name)
				if token == "" {
					// First load — create a new progress token.
					tokenID++
					token = fmt.Sprintf("tsgo-loading-%d", tokenID)
					begun = false
					p.sendProgressCreate(token)
				}
				msg := p.displayName(ev.name)
				if !begun {
					begun = true
					p.sendProgress(token, &lsproto.WorkDoneProgressBegin{
						Title:   "Loading",
						Message: &msg,
					})
				} else {
					p.sendProgress(token, &lsproto.WorkDoneProgressReport{
						Message: &msg,
					})
				}
			} else {
				loading.Delete(ev.name)
				if token == "" {
					continue
				}
				if loading.Size() == 0 {
					// Last project finished — end the progress.
					p.sendProgress(token, &lsproto.WorkDoneProgressEnd{})
					token = ""
				} else {
					// Show the oldest still-loading project.
					first := firstValue(loading.Values())
					msg := p.displayName(first)
					p.sendProgress(token, &lsproto.WorkDoneProgressReport{
						Message: &msg,
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

// displayName returns a short display string for a project name,
// making it relative to the workspace root and truncating if needed.
func (p *projectLoadingProgress) displayName(projectName string) string {
	projectName = tspath.ConvertToRelativePath(projectName, tspath.ComparePathsOptions{
		CurrentDirectory: p.server.cwd,
	})
	const maxLen = 60
	if len(projectName) > maxLen {
		projectName = "..." + projectName[len(projectName)-maxLen:]
	}
	return projectName
}

func firstValue[T any](seq func(yield func(T) bool)) T {
	for v := range seq {
		return v
	}
	panic("empty sequence")
}
