// Package checker re-exports the internal type checker API for external
// consumers of this fork. See package github.com/microsoft/typescript-go/shim
// for the governing policy.
package checker

import (
	"github.com/microsoft/typescript-go/internal/checker"
)

type (
	Checker = checker.Checker
	Type    = checker.Type
)
