package pnp

func InitPnpApi(fs PnpApiFS, filePath string) *PnpApi {
	pnpApi := &PnpApi{fs: fs, url: filePath}

	manifestData, err := pnpApi.findClosestPnpManifest()
	if err == nil {
		pnpApi.manifest = manifestData
		return pnpApi
	}

	return nil
}
