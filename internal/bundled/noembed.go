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

// argv0Dir returns the directory the executable was invoked from
// (os.Args[0], resolved via PATH when bare, then symlink-walked per
// component). It is the fallback when /proc/self/exe names the executable
// somewhere its neighbours are not visible from — e.g. under PRoot's
// link2symlink extension (Termux proot-distro on Android), which executes
// hardlinked binaries out of a hidden symlink store that only the kernel
// sees. Returns "" when no usable directory can be derived.
func argv0Dir() string {
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
	return tspath.GetDirectoryPath(tspath.NormalizeSlashes(abs))
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

	if fallback := argv0Dir(); fallback != "" && fallback != dir && hasLibDTS(fallback) {
		return fallback
	}

	panic(fmt.Sprintf("bundled: %v does not exist; this executable may be misplaced", tspath.CombinePaths(dir, "lib.d.ts")))
})

func IsBundled(path string) bool {
	return false
}
