//go:build !windows

package osvfs

// isSymlinkOrJunction always returns false on non-Windows platforms.
// On Unix-like systems, symlinks are already properly detected by the
// fs.ModeSymlink bit in the directory entry type, so this check is not needed.
func isSymlinkOrJunction(path string) bool {
	return false
}
