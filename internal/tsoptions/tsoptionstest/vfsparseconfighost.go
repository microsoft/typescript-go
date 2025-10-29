package tsoptionstest

import (
	"github.com/microsoft/typescript-go/internal/pnp"
	"github.com/microsoft/typescript-go/internal/tsoptions"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs"
	"github.com/microsoft/typescript-go/internal/vfs/pnpvfs"
	"github.com/microsoft/typescript-go/internal/vfs/vfstest"
)

func fixRoot(path string) string {
	rootLength := tspath.GetRootLength(path)
	if rootLength == 0 {
		return path
	}
	if len(path) == rootLength {
		return "."
	}
	return path[rootLength:]
}

type VfsParseConfigHost struct {
	Vfs              vfs.FS
	CurrentDirectory string
	pnpApi           *pnp.PnpApi
}

var _ tsoptions.ParseConfigHost = (*VfsParseConfigHost)(nil)

func (h *VfsParseConfigHost) FS() vfs.FS {
	return h.Vfs
}

func (h *VfsParseConfigHost) GetCurrentDirectory() string {
	return h.CurrentDirectory
}

func (h *VfsParseConfigHost) PnpApi() *pnp.PnpApi {
	return h.pnpApi
}

func NewVFSParseConfigHost(files map[string]string, currentDirectory string, useCaseSensitiveFileNames bool) *VfsParseConfigHost {
	var fs vfs.FS = vfstest.FromMap(files, useCaseSensitiveFileNames)
	pnpApi := pnp.InitPnpApi(fs, tspath.NormalizePath(currentDirectory))
	if pnpApi != nil {
		fs = pnpvfs.From(fs)
	}

	return &VfsParseConfigHost{
		Vfs:              fs,
		CurrentDirectory: currentDirectory,
		pnpApi:           nil,
	}
}
