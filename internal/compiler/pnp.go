package compiler

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/gun-yu/pnp-go/pkg"
)

func TryGetPnpResolutionConfig(path string) *pkg.ResolutionConfig {
	pnpManifestPath, err := findNearestPNPPath(path)
	if err != nil {
		return nil
	}
	pnpManifest, err := pkg.LoadPNPManifest(pnpManifestPath)
	if err != nil {
		return nil
	}

	return &pkg.ResolutionConfig{
		Host: pkg.ResolutionHost{
			FindPNPManifest: func(_ string) (*pkg.Manifest, error) {
				return &pnpManifest, nil
			},
		},
	}
}

func findNearestPNPPath(start string) (string, error) {
	dir := start
	if fi, err := os.Stat(start); err == nil {
		if !fi.IsDir() {
			dir = filepath.Dir(start)
		}
	} else {
		dir = filepath.Dir(start)
	}

	for {
		for _, name := range []string{".pnp.data.json", ".pnp.cjs", ".pnp.js"} {
			candidate := filepath.Join(dir, name)
			if fi, err := os.Stat(candidate); err == nil && !fi.IsDir() {
				return candidate, nil
			}
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", fs.ErrNotExist
}
