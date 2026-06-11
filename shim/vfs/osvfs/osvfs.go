// Package osvfs re-exports the internal OS-backed file system for external
// consumers of this fork. See package github.com/microsoft/typescript-go/shim
// for the governing policy.
package osvfs

import (
	"github.com/microsoft/typescript-go/internal/vfs/osvfs"
)

var FS = osvfs.FS
