package lspservertests

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/microsoft/typescript-go/internal/testutil/baseline"
)

type lspServerTest struct {
	subscenario string
	files       func() map[string]any
	test        func(server *testServer)
}

func (test *lspServerTest) run(t *testing.T, scenario string) {
	t.Helper()
	t.Run(scenario+"/"+test.subscenario, func(t *testing.T) {
		t.Parallel()
		server := newTestServer(t, test.files())
		test.test(server)
		baseline.Run(t, strings.ReplaceAll(test.subscenario, " ", "-")+".js", server.baseline.String(), baseline.Options{Subfolder: filepath.Join("lspservertests", scenario)})
	})
}
