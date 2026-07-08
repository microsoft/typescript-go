package pnp

import (
	"github.com/microsoft/typescript-go/internal/tspath"
)

// Manifest filenames, in the discovery order Yarn and esbuild use. Yarn writes
// .pnp.cjs always; .pnp.data.json additionally when pnpEnableInlining is false.
var manifestNames = []string{".pnp.data.json", ".pnp.cjs", ".pnp.js"}

// Find walks up from startDir (an absolute, normalized directory) looking for a
// PnP manifest, and returns the absolute path of the manifest file to load
// (preferring a .pnp.data.json sidecar, else the .pnp.cjs the data is inlined
// into), or "" if there is no PnP manifest above startDir.
func Find(startDir string, fileExists func(string) bool) string {
	manifest, _ := tspath.ForEachAncestorDirectory(tspath.NormalizePath(startDir), func(dir string) (string, bool) {
		for _, name := range manifestNames {
			candidate := tspath.CombinePaths(dir, name)
			if fileExists(candidate) {
				return candidate, true
			}
		}
		return "", false
	})
	return manifest
}
