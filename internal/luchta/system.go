package luchta

import (
	"io"
	"time"

	"github.com/microsoft/typescript-go/internal/pnp"
	"github.com/microsoft/typescript-go/internal/vfs"
)

// runSystem satisfies internal/execute/tsc.System and internal/tsoptions.ParseConfigHost
// for a single package compilation rooted at cwd. Writer routes diagnostic text to the
// provided io.Writer (never to stdout).
type runSystem struct {
	cwd    string
	fs     vfs.FS
	libs   string
	out    io.Writer
	start  time.Time
	pnpApi *pnp.PnpApi
}

func newRunSystem(cwd string, fsys vfs.FS, libraryPath string, w io.Writer, pnpApi *pnp.PnpApi) *runSystem {
	return &runSystem{cwd: cwd, fs: fsys, libs: libraryPath, out: w, start: time.Now(), pnpApi: pnpApi}
}

func (s *runSystem) Writer() io.Writer                         { return s.out }
func (s *runSystem) FS() vfs.FS                                { return s.fs }
func (s *runSystem) DefaultLibraryPath() string                { return s.libs }
func (s *runSystem) GetCurrentDirectory() string               { return s.cwd }
func (s *runSystem) WriteOutputIsTTY() bool                    { return false }
func (s *runSystem) GetWidthOfTerminal() int                   { return 0 }
func (s *runSystem) GetEnvironmentVariable(name string) string { return "" }
func (s *runSystem) Now() time.Time                            { return time.Now() }
func (s *runSystem) SinceStart() time.Duration                 { return time.Since(s.start) }
func (s *runSystem) PnpApi() *pnp.PnpApi                       { return s.pnpApi }
