package watchmanager

import (
	"slices"

	"github.com/microsoft/typescript-go/internal/ast"
)

// IsModuleResolutionError returns true if the diagnostic indicates
// a failed module resolution (TS2307, TS2732, TS2792, TS2882).
func IsModuleResolutionError(d *ast.Diagnostic) bool {
	switch d.Code() {
	case 2307, // Cannot find module '{0}' or its corresponding type declarations.
		2732, // Cannot find module '{0}'. Consider using '--resolveJsonModule'...
		2792, // Cannot find module '{0}'. Did you mean to set the 'moduleResolution' option...
		2882: // Cannot find module or type declarations for side-effect import of '{0}'.
		return true
	}
	return false
}

// HasModuleResolutionErrors returns true if any of the diagnostics indicate
// a failed module resolution.
func HasModuleResolutionErrors(diagnostics []*ast.Diagnostic) bool {
	return slices.ContainsFunc(diagnostics, func(d *ast.Diagnostic) bool {
		return IsModuleResolutionError(d)
	})
}
