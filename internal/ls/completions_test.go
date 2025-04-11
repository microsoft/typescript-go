package ls_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/ls"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/project"
	"gotest.tools/v3/assert"
)

func TestCompletions(t *testing.T) {
	files := map[string]string{
		"index.ts": "",
	}
	ls := createLanguageService("/", files)
	var context *lsproto.CompletionContext
	var capabilities *lsproto.CompletionClientCapabilities
	completionList := ls.ProvideCompletion("index.ts", 0, context, capabilities)
	assert.Assert(t, completionList != nil)
}

func createLanguageService(cd string, files map[string]string) *ls.LanguageService {
	// !!! TODO: replace with service_test.go's `setup`
	projectServiceHost := newProjectServiceHost(files)
	projectService := project.NewService(projectServiceHost, project.ServiceOptions{})
	compilerOptions := &core.CompilerOptions{}
	project := project.NewInferredProject(compilerOptions, cd, "/", projectService)
	return project.LanguageService()
}

func newProjectServiceHost(files map[string]string) project.ServiceHost {
	// !!! TODO: import from service_test.go
	return nil
}
