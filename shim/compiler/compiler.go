// Package compiler re-exports the internal compiler API (Program and
// CompilerHost) for external consumers of this fork. See package
// github.com/microsoft/typescript-go/shim for the governing policy.
package compiler

import (
	"github.com/microsoft/typescript-go/internal/compiler"
)

type (
	Program        = compiler.Program
	ProgramOptions = compiler.ProgramOptions
	CompilerHost   = compiler.CompilerHost
)

var (
	NewProgram              = compiler.NewProgram
	NewCompilerHost         = compiler.NewCompilerHost
	NewCachedFSCompilerHost = compiler.NewCachedFSCompilerHost
)
