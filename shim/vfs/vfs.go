// Package vfs re-exports the internal virtual file system API for external
// consumers of this fork. See package github.com/microsoft/typescript-go/shim
// for the governing policy.
package vfs

import (
	"github.com/microsoft/typescript-go/internal/vfs"
)

type FS = vfs.FS
