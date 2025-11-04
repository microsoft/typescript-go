//go:build !windows

package osvfs

// On Unix-like systems, symlinks are already properly detected by the
// fs.ModeSymlink bit in the directory entry type, so this check is not needed.
var isSymlinkOrJunction func(path string) bool
