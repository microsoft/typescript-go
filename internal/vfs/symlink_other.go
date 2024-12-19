//go:build !windows

package vfs

import "os"

func isJunction(fi os.FileInfo) bool {
	return false
}
