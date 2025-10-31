package lspservertests

import (
	"strings"
	"testing"

	"github.com/microsoft/typescript-go/internal/testutil/baseline"
	"github.com/microsoft/typescript-go/internal/tspath"
)

type lspServerTest struct {
	subScenario string
	files       func() map[string]any
	test        func(server *testServer)
}

func (test *lspServerTest) run(t *testing.T, scenario string) {
	t.Helper()
	t.Run(scenario+"/"+test.subScenario, func(t *testing.T) {
		t.Parallel()
		server := newTestServer(t, test.files())
		test.test(server)
		baseline.Run(t, strings.ReplaceAll(test.subScenario, " ", "-")+".js", server.baseline.String(), baseline.Options{Subfolder: tspath.CombinePaths("lspservertests", scenario)})
	})
}
