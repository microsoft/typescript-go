package contentmappertest_test

import (
	"context"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/microsoft/typescript-go/internal/contentmapper"
	"github.com/microsoft/typescript-go/internal/contentmapperhost"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/spanmap"
	"github.com/microsoft/typescript-go/internal/testutil/contentmappertest"
	"gotest.tools/v3/assert"
)

// helperEnv, when set, makes the test binary act as the mapper subprocess instead of running tests. This
// lets the out-of-process test spawn a real subprocess (itself) that speaks the mapper protocol over
// stdio, exercising the same handler code that the in-process spawner runs over a pipe.
const helperEnv = "TSGO_CONTENT_MAPPER_HELPER"

func TestMain(m *testing.M) {
	if os.Getenv(helperEnv) == "1" {
		_ = contentmappertest.Serve(context.Background(), stdio{})
		os.Exit(0)
	}
	os.Exit(m.Run())
}

// stdio adapts the process's stdin/stdout to an io.ReadWriteCloser for the mapper server.
type stdio struct{}

func (stdio) Read(p []byte) (int, error)  { return os.Stdin.Read(p) }
func (stdio) Write(p []byte) (int, error) { return os.Stdout.Write(p) }
func (stdio) Close() error                { return nil }

func testMapper() *contentmapper.Mapper {
	return &contentmapper.Mapper{
		Definition: contentmapper.Definition{
			Package:    contentmappertest.PackageName,
			Extensions: []string{".box"},
		},
		Manifest: contentmapper.Manifest{
			Name:    contentmappertest.PackageName,
			Version: "1.0.0",
			Exec:    []string{contentmappertest.ExecName},
		},
		PackageDirectory: "/node_modules/" + contentmappertest.PackageName,
	}
}

func transformRequest() contentmapperhost.Request {
	return contentmapperhost.Request{
		FileName:        "/app.box",
		Content:         "export const version = #{target};\n",
		ConfigFileName:  "/tsconfig.json",
		CompilerOptions: &core.CompilerOptions{Target: core.ScriptTargetES2020},
	}
}

// TestInProcessSpanKinds drives the mapper in-process and verifies the transform produces a span map that
// exercises all three span kinds: the synthesized preamble, the verbatim body, and the atom substitution
// of a compiler-option token.
func TestInProcessSpanKinds(t *testing.T) {
	t.Parallel()
	host := contentmapperhost.New(t.Context(), contentmappertest.NewSpawner())
	defer host.Close()

	result, err := host.Transform(testMapper(), transformRequest())
	assert.NilError(t, err)
	assert.Equal(t, result.ScriptKind, core.ScriptKindTS)
	// The #{target} token was replaced by the es2020 target value (7).
	assert.Assert(t, strings.Contains(result.Text, "export const version = 7;"), "got %q", result.Text)

	text := result.Text

	// Synthesized: __VERSION appears only in the injected preamble, which has no original counterpart.
	synthStart := strings.Index(text, "__VERSION")
	assert.Assert(t, synthStart >= 0)
	_, synthFidelity := result.Mappings.MapSpan(core.NewTextRange(synthStart, synthStart+len("__VERSION")))
	assert.Equal(t, synthFidelity, spanmap.FidelityNone)

	// Verbatim: "export const version" is copied through, so its generated span maps exactly onto the
	// identical span in the original.
	verbatimStart := strings.Index(text, "export const version")
	assert.Assert(t, verbatimStart >= 0)
	verbatimRange, verbatimFidelity := result.Mappings.MapSpan(core.NewTextRange(verbatimStart, verbatimStart+len("export")))
	assert.Equal(t, verbatimFidelity, spanmap.FidelityExact)
	original := transformRequest().Content
	assert.Equal(t, original[verbatimRange.Pos():verbatimRange.End()], "export")

	// Atom: the substituted "7" maps as a whole back to the original #{target} token span.
	atomStart := strings.Index(text, "= 7;") + len("= ")
	atomRange, atomFidelity := result.Mappings.MapSpan(core.NewTextRange(atomStart, atomStart+len("7")))
	assert.Equal(t, atomFidelity, spanmap.FidelityAtom)
	assert.Equal(t, original[atomRange.Pos():atomRange.End()], "#{target}")
}

// TestOutOfProcess exercises the real out-of-process IPC path: it spawns the test binary as a mapper
// subprocess and drives it over stdio through the production content mapper host.
func TestOutOfProcess(t *testing.T) {
	t.Parallel()
	host := contentmapperhost.New(t.Context(), execSpawner{})
	defer host.Close()

	result, err := host.Transform(testMapper(), transformRequest())
	assert.NilError(t, err)
	assert.Equal(t, result.ScriptKind, core.ScriptKindTS)
	assert.Assert(t, strings.Contains(result.Text, "export const version = 7;"), "got %q", result.Text)
	assert.Assert(t, result.Mappings != nil)
}

// execSpawner spawns the test binary itself as the mapper subprocess (guarded by helperEnv), so the test
// talks to a genuinely separate process over real pipes.
type execSpawner struct{}

func (execSpawner) Spawn(command []string, dir string) (io.ReadWriteCloser, error) {
	cmd := exec.Command(os.Args[0])
	cmd.Env = append(os.Environ(), helperEnv+"=1")
	cmd.Stderr = os.Stderr
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	return &process{cmd: cmd, stdin: stdin, stdout: stdout}, nil
}

// process adapts a spawned subprocess's stdio to an io.ReadWriteCloser: reads come from its stdout, writes
// go to its stdin, and Close tears the process down.
type process struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout io.ReadCloser
}

func (p *process) Read(b []byte) (int, error)  { return p.stdout.Read(b) }
func (p *process) Write(b []byte) (int, error) { return p.stdin.Write(b) }

func (p *process) Close() error {
	_ = p.stdin.Close()
	_ = p.cmd.Process.Kill()
	_ = p.cmd.Wait()
	return nil
}
