//go:build noembed

package bundled

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"testing"

	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs"
	"github.com/microsoft/typescript-go/internal/vfs/osvfs"
)

const embedded = false

func wrapFS(fs vfs.FS) vfs.FS {
	return fs
}

var executableDir = sync.OnceValue(func() string {
	exe, err := os.Executable()
	if err != nil {
		panic(fmt.Sprintf("bundled: failed to get executable path: %v", err))
	}
	exe = tspath.NormalizeSlashes(exe)
	exe = osvfs.FS().Realpath(exe)
	return tspath.GetDirectoryPath(exe)
})

// argv0Path returns the invocation path of the executable (os.Args[0],
// resolved via PATH when bare, then symlink-walked), as a fallback for when
// /proc/self/exe reports a path this process cannot see (e.g. under PRoot's
// link2symlink). Returns "" when no usable path can be derived.
func argv0Path() string {
	if len(os.Args) == 0 || os.Args[0] == "" {
		return ""
	}
	argv0, err := exec.LookPath(os.Args[0])
	if err != nil && !errors.Is(err, exec.ErrDot) {
		return ""
	}
	abs, err := filepath.Abs(argv0)
	if err != nil {
		return ""
	}
	if walked, err := filepath.EvalSymlinks(abs); err == nil {
		abs = walked
	}
	return tspath.NormalizeSlashes(abs)
}

// sameAsExecutable reports whether path provably names the same file the
// kernel reports as the running executable.
func sameAsExecutable(path string) bool {
	exe, err := os.Executable()
	if err != nil {
		return false
	}
	exeInfo, err := os.Stat(exe)
	if err != nil {
		return false
	}
	info, err := os.Stat(path)
	return err == nil && os.SameFile(exeInfo, info)
}

func hasLibDTS(dir string) bool {
	return osvfs.FS().Stat(tspath.CombinePaths(dir, "lib.d.ts")) != nil
}

var libPath = sync.OnceValue(func() string {
	if testing.Testing() {
		return TestingLibPath()
	}
	dir := executableDir()
	if hasLibDTS(dir) {
		return dir
	}

	// Only trust argv[0] when it names the very file the kernel says is
	// running; a different file could be another installation with
	// mismatched libs.
	fallback := ""
	if argv0 := argv0Path(); argv0 != "" && sameAsExecutable(argv0) {
		fallback = tspath.GetDirectoryPath(argv0)
		if fallback != dir && hasLibDTS(fallback) {
			return fallback
		}
	}

	msg := fmt.Sprintf("bundled: %v does not exist", tspath.CombinePaths(dir, "lib.d.ts"))
	if fallback != "" && fallback != dir {
		msg += fmt.Sprintf(" (also checked %v)", tspath.CombinePaths(fallback, "lib.d.ts"))
	}
	panic(msg + "; this executable may be misplaced")
})

func IsBundled(path string) bool {
	return false
}
