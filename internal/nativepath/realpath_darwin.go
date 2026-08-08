package nativepath

import (
	"os"
	"path/filepath"
	"sync"
	"unsafe"

	"golang.org/x/sys/unix"
)

// On macOS, we use open + fcntl(F_GETPATH) to resolve the canonical path in
// O(1) syscalls instead of Go's filepath.EvalSymlinks which does an lstat per
// path component — O(depth).
//
// How it works:
//   - open(path, O_EVTONLY|O_NONBLOCK|O_CLOEXEC) follows all symlinks and gives
//     us a lightweight fd. O_EVTONLY is macOS's event-only descriptor — it
//     doesn't require read permission (similar to Linux's O_PATH) but still
//     references the vnode. O_NONBLOCK prevents blocking on FIFOs.
//   - fcntl(fd, F_GETPATH, buf) asks the kernel for the canonical path of the
//     open file descriptor, written into a MAXPATHLEN buffer.
//
// F_GETPATH answers from the vnode name cache, so for a regular file with
// nlink > 1 it may return any hardlink sibling's name — nondeterministically
// under concurrency (pnpm and Nix stores hardlink-dedupe identical files).
// For that case we resolve the parent directory instead (directories cannot
// be hardlinked) and re-attach the final component.
//
// unix.FcntlInt takes an int arg, so call it through a uintptr-escaping wrapper
// to keep the buffer pointer valid until fcntl returns.

var hasFGetPath = sync.OnceValue(func() bool {
	// Verify that F_GETPATH is supported by this kernel version.
	var buf [unix.PathMax]byte
	fd, err := unix.Open(".", unix.O_EVTONLY|unix.O_NONBLOCK|unix.O_CLOEXEC, 0)
	if err != nil {
		return false
	}
	defer unix.Close(fd)
	_, err = fcntlGetPath(fd, &buf)
	return err == nil
})

func fcntlGetPath(fd int, buf *[unix.PathMax]byte) (int, error) {
	return ignoringEINTR(func() (int, error) {
		return fcntlGetPathPtr(uintptr(fd), uintptr(unsafe.Pointer(&buf[0])))
	})
}

//go:uintptrescapes
func fcntlGetPathPtr(fd uintptr, buf uintptr) (int, error) {
	return unix.FcntlInt(fd, unix.F_GETPATH, int(buf))
}

// maxSymlinkDepth mirrors the kernel's symlink loop limit (MAXSYMLINKS).
const maxSymlinkDepth = 40

func Realpath(path string) (string, error) {
	if !hasFGetPath() {
		return filepath.EvalSymlinks(path)
	}
	return realpathFast(path, 0)
}

func realpathFast(path string, depth int) (string, error) {
	fd, err := unix.Open(path, unix.O_EVTONLY|unix.O_NONBLOCK|unix.O_CLOEXEC, 0)
	if err != nil {
		return "", err
	}

	var st unix.Stat_t
	if _, err := ignoringEINTR(func() (struct{}, error) {
		return struct{}{}, unix.Fstat(fd, &st)
	}); err != nil {
		unix.Close(fd)
		return "", err
	}

	if st.Mode&unix.S_IFMT == unix.S_IFDIR || st.Nlink == 1 {
		defer unix.Close(fd)
		var buf [unix.PathMax]byte
		if _, err := fcntlGetPath(fd, &buf); err != nil {
			return "", err
		}
		return unix.ByteSliceToString(buf[:]), nil
	}
	unix.Close(fd)

	// hardlinked regular file: F_GETPATH may name any of its links (see above)
	if depth > maxSymlinkDepth {
		return "", unix.ELOOP
	}
	dir, base := filepath.Dir(path), filepath.Base(path)
	if dir == path || base == "." || base == ".." || base == string(filepath.Separator) {
		// Cannot split the path further; defer to the stdlib resolver.
		return filepath.EvalSymlinks(path)
	}
	resolvedDir, err := realpathFast(dir, depth)
	if err != nil {
		return "", err
	}
	joined := filepath.Join(resolvedDir, base)
	fi, err := os.Lstat(joined)
	if err != nil {
		return "", err
	}
	if fi.Mode()&os.ModeSymlink == 0 {
		return joined, nil
	}
	target, err := os.Readlink(joined)
	if err != nil {
		return "", err
	}
	if !filepath.IsAbs(target) {
		target = filepath.Join(resolvedDir, target)
	}
	return realpathFast(target, depth+1)
}
