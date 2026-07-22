package nativepath

import (
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"golang.org/x/sys/unix"
)

// On Linux, we use the O_PATH + /proc/self/fd trick to resolve the canonical
// path in O(1) syscalls (open + readlink + close) instead of Go's
// filepath.EvalSymlinks which does an lstat per path component — O(depth).
//
// This is the approach libuv/Node.js could use, though libuv currently just
// calls C realpath(3) which itself does a readlink per component. On the Go
// side, the per-component approach is even more expensive because each
// os.Lstat call involves goroutine scheduling overhead (entersyscall /
// exitsyscall).
//
// How it works:
//   - open(path, O_PATH|O_CLOEXEC) gives us a lightweight fd that follows all
//     symlinks to the final target. O_PATH requires only search permission on
//     directories (same as lstat), and works for both files and directories.
//   - readlink("/proc/self/fd/<fd>") returns the fully resolved canonical path
//     that the kernel computed during the open.
//
// Falls back to filepath.EvalSymlinks if /proc is not available (e.g. containers
// or chroots without procfs mounted).
//
// The kernel's answer is validated with one invariant: canonicalization may
// only rename the final path component when that component is a symlink. If
// the name changed but lstat reports a non-symlink, the filesystem view is
// being virtualized behind the process's back — e.g. PRoot's link2symlink
// extension (the default under Termux proot-distro on Android), which
// emulates hardlinks with symlinks into a hidden store and masks them from
// lstat/readlink, but cannot mask names the kernel hands back through /proc.
// In that case we fall back to the per-component walk, which sees the same
// masked view as glibc's realpath(3) and returns a path that is actually
// usable by this process. The check costs no extra syscalls on the common
// path (a string comparison), and one lstat when the final component was
// renamed.

const _procSelfFD = "/proc/self/fd/"

var hasProcSelfFD = sync.OnceValue(func() bool {
	var stat unix.Stat_t
	return unix.Stat(_procSelfFD, &stat) == nil
})

func Realpath(path string) (string, error) {
	if !hasProcSelfFD() {
		return filepath.EvalSymlinks(path)
	}

	resolved, err := procRealpath(path)
	if err != nil {
		return "", err
	}

	if filepath.Base(resolved) != filepath.Base(path) {
		if info, err := os.Lstat(path); err == nil && info.Mode()&os.ModeSymlink == 0 {
			if walked, err := filepath.EvalSymlinks(path); err == nil {
				return walked, nil
			}
		}
	}

	return resolved, nil
}

func procRealpath(path string) (string, error) {
	fd, err := ignoringEINTR(func() (int, error) {
		return unix.Open(path, unix.O_CLOEXEC|unix.O_PATH, 0)
	})
	if err != nil {
		return "", &os.PathError{Op: "open", Path: path, Err: err}
	}
	defer unix.Close(fd)

	var procBuf [len(_procSelfFD) + 20]byte // 20 digits is enough for any int64 fd
	n := copy(procBuf[:], _procSelfFD)
	n += copy(procBuf[n:], strconv.Itoa(fd))
	procPath := string(procBuf[:n])

	buf := make([]byte, 256)
	for {
		nn, err := ignoringEINTR(func() (int, error) {
			return unix.Readlink(procPath, buf)
		})
		if err != nil {
			return "", &os.PathError{Op: "readlink", Path: path, Err: err}
		}
		if nn < len(buf) {
			return string(buf[:nn]), nil
		}
		buf = make([]byte, len(buf)*2)
	}
}
