package osvfs

import (
	"syscall"
	"unsafe"
)

// isSymlinkOrJunction checks if the given path is a symlink or junction point
// on Windows by checking the FILE_ATTRIBUTE_REPARSE_POINT attribute.
func isSymlinkOrJunction(path string) bool {
	if len(path) >= 248 {
		path = `\\?\` + path
	}

	pathUTF16, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return false
	}

	var data syscall.Win32FileAttributeData
	err = syscall.GetFileAttributesEx(
		pathUTF16,
		syscall.GetFileExInfoStandard,
		(*byte)(unsafe.Pointer(&data)),
	)
	if err != nil {
		return false
	}

	return data.FileAttributes&syscall.FILE_ATTRIBUTE_REPARSE_POINT != 0
}
