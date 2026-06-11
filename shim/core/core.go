// Package core re-exports internal core types for external consumers of this
// fork. See package github.com/microsoft/typescript-go/shim for the governing
// policy.
package core

import (
	"github.com/microsoft/typescript-go/internal/core"
)

type (
	CompilerOptions = core.CompilerOptions
	Tristate        = core.Tristate
	ScriptKind      = core.ScriptKind
	ScriptTarget    = core.ScriptTarget
)

const (
	TSUnknown = core.TSUnknown
	TSFalse   = core.TSFalse
	TSTrue    = core.TSTrue
)
