package vfs

import (
	"os"
	"syscall"
)

func isJunction(fi os.FileInfo) bool {
	sys := fi.Sys().(*syscall.Win32FileAttributeData)
	return sys.FileAttributes&syscall.FILE_ATTRIBUTE_REPARSE_POINT != 0 && sys.FileAttributes&syscall.FILE_ATTRIBUTE_DIRECTORY != 0
}
