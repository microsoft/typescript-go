package osvfs

func isSymlink(path string) bool {
	return isReparsePoint(path)
}
