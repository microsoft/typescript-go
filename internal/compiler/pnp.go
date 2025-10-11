package compiler

import (
	"io/fs"
	"os"
	"path/filepath"

	module "github.com/microsoft/typescript-go/internal/module/pnp"
)

func TryGetPnpResolutionConfig(path string) *module.ResolutionConfig {
	pnpManifestPath, err := findNearestPNPPath(path)
	if err != nil {
		return nil
	}
	pnpManifest, err := module.LoadPNPManifest(pnpManifestPath)
	if err != nil {
		return nil
	}

	return &module.ResolutionConfig{
		Host: module.ResolutionHost{
			FindPNPManifest: func(_ string) (*module.Manifest, error) {
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
