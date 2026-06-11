// Package bundled re-exports access to the embedded lib.d.ts files for
// external consumers of this fork. See package
// github.com/microsoft/typescript-go/shim for the governing policy.
package bundled

import (
	"github.com/microsoft/typescript-go/internal/bundled"
)

const Embedded = bundled.Embedded

var (
	WrapFS  = bundled.WrapFS
	LibPath = bundled.LibPath
)
