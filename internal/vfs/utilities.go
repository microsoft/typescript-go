package vfs

import (
	"strings"
)

// An "includes" path "foo" is implicitly a glob "foo/** /*" (without the space) if its last component has no extension,
// and does not contain any glob characters itself.
func IsImplicitGlob(lastPathComponent string) bool {
	return !strings.ContainsAny(lastPathComponent, ".*?")
}
