// Package tsoptions re-exports the internal tsconfig parsing API for external
// consumers of this fork. See package github.com/microsoft/typescript-go/shim
// for the governing policy.
package tsoptions

import (
	"github.com/microsoft/typescript-go/internal/tsoptions"
)

type (
	ParsedCommandLine   = tsoptions.ParsedCommandLine
	ParseConfigHost     = tsoptions.ParseConfigHost
	ExtendedConfigCache = tsoptions.ExtendedConfigCache
)

var GetParsedCommandLineOfConfigFile = tsoptions.GetParsedCommandLineOfConfigFile
