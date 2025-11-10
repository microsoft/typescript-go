//go:build !windows

package osvfs

// On Unix-like systems, symlinks are already properly detected by the
// fs.ModeSymlink bit in the directory entry type, so this check is not needed.
var isReparsePoint func(path string) bool
