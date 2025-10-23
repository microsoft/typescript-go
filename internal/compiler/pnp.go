package compiler

import (
	"github.com/microsoft/typescript-go/internal/module/pnp"
	"github.com/microsoft/typescript-go/internal/vfs"
)

func TryGetPnpResolutionConfig(path string, fs vfs.FS) *pnp.ResolutionConfig {
	pnpManifest, err := pnp.FindPNPManifest(path, fs)
	if err != nil {
		return nil
	}

	return &pnp.ResolutionConfig{
		Host: pnp.PNPResolutionHost{
			FindPNPManifest: func(_ string) (*pnp.Manifest, error) {
				return pnpManifest, nil
			},
		},
	}
}
