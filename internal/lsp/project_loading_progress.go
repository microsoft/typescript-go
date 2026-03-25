package lsp

import (
	"context"
	"fmt"
	"sync"

	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/tspath"
)

// projectLoadingProgress manages a single LSP WorkDoneProgress token that
// stays active while any project is loading. Multiple concurrent loads are
// tracked; the message shows the most recently started project, and the
// indicator disappears only when all loads finish.
type projectLoadingProgress struct {
	server  *Server
	mu      sync.Mutex
	loading collections.OrderedSet[string]
	token   string // empty if no progress is active
	tokenID int
}

func newProjectLoadingProgress(server *Server) *projectLoadingProgress {
	return &projectLoadingProgress{
		server: server,
	}
}

func (p *projectLoadingProgress) start(ctx context.Context, projectName string) {
	var newToken string

	p.mu.Lock()
	p.loading.Add(projectName)
	isFirst := p.token == ""
	if isFirst {
		p.tokenID++
		p.token = fmt.Sprintf("tsgo-loading-%d", p.tokenID)
	}
	newToken = p.token
	p.mu.Unlock()

	msg := p.displayName(projectName)
	if isFirst {
		_, _ = sendClientRequest(ctx, p.server, lsproto.WindowWorkDoneProgressCreateInfo, &lsproto.WorkDoneProgressCreateParams{
			Token: lsproto.IntegerOrString{String: &newToken},
		})

		_ = sendNotification(p.server, lsproto.ProgressInfo, &lsproto.ProgressParams{
			Token: lsproto.IntegerOrString{String: &newToken},
			Value: &lsproto.WorkDoneProgressBegin{
				Title:   "Loading",
				Message: &msg,
			},
		})
	} else {
		_ = sendNotification(p.server, lsproto.ProgressInfo, &lsproto.ProgressParams{
			Token: lsproto.IntegerOrString{String: &newToken},
			Value: &lsproto.WorkDoneProgressReport{
				Message: &msg,
			},
		})
	}
}

func (p *projectLoadingProgress) finish(ctx context.Context, projectName string) {
	p.mu.Lock()
	p.loading.Delete(projectName)

	tokenStr := p.token
	if tokenStr == "" {
		p.mu.Unlock()
		return
	}

	done := p.loading.Size() == 0
	var first string
	if done {
		p.token = ""
	} else {
		first = p.firstLoading()
	}
	p.mu.Unlock()

	if done {
		_ = sendNotification(p.server, lsproto.ProgressInfo, &lsproto.ProgressParams{
			Token: lsproto.IntegerOrString{String: &tokenStr},
			Value: &lsproto.WorkDoneProgressEnd{},
		})
	} else {
		msg := p.displayName(first)
		_ = sendNotification(p.server, lsproto.ProgressInfo, &lsproto.ProgressParams{
			Token: lsproto.IntegerOrString{String: &tokenStr},
			Value: &lsproto.WorkDoneProgressReport{
				Message: &msg,
			},
		})
	}
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

// firstLoading returns the oldest project still in the loading set.
// Must be called with p.mu held and p.loading.Size() > 0.
func (p *projectLoadingProgress) firstLoading() string {
	for name := range p.loading.Values() {
		return name
	}
	panic("unreachable")
}
