package pnpvfs

import (
	"archive/zip"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs"
	"github.com/microsoft/typescript-go/internal/vfs/iovfs"
)

type pnpFS struct {
	fs                  vfs.FS
	cachedZipReadersMap map[string]*zip.ReadCloser
	cacheReaderMutex    sync.Mutex
}

var _ vfs.FS = (*pnpFS)(nil)

func From(fs vfs.FS) *pnpFS {
	pnpFS := &pnpFS{
		fs:                  fs,
		cachedZipReadersMap: make(map[string]*zip.ReadCloser),
	}

	return pnpFS
}

func (pnpFS *pnpFS) DirectoryExists(path string) bool {
	path, _, _ = resolveVirtual(path)

	if strings.HasSuffix(path, ".zip") {
		return pnpFS.fs.FileExists(path)
	}

	fs, formattedPath, _ := getMatchingFS(pnpFS, path)

	return fs.DirectoryExists(formattedPath)
}

func (pnpFS *pnpFS) FileExists(path string) bool {
	path, _, _ = resolveVirtual(path)

	if strings.HasSuffix(path, ".zip") {
		return pnpFS.fs.FileExists(path)
	}

	fs, formattedPath, _ := getMatchingFS(pnpFS, path)
	return fs.FileExists(formattedPath)
}

func (pnpFS *pnpFS) GetAccessibleEntries(path string) vfs.Entries {
	path, hash, basePath := resolveVirtual(path)

	fs, formattedPath, zipPath := getMatchingFS(pnpFS, path)
	entries := fs.GetAccessibleEntries(formattedPath)

	for i, dir := range entries.Directories {
		fullPath := tspath.CombinePaths(zipPath+formattedPath, dir)
		entries.Directories[i] = makeVirtualPath(basePath, hash, fullPath)
	}

	for i, file := range entries.Files {
		fullPath := tspath.CombinePaths(zipPath+formattedPath, file)
		entries.Files[i] = makeVirtualPath(basePath, hash, fullPath)
	}

	return entries
}

func (pnpFS *pnpFS) ReadFile(path string) (contents string, ok bool) {
	path, _, _ = resolveVirtual(path)

	fs, formattedPath, _ := getMatchingFS(pnpFS, path)
	return fs.ReadFile(formattedPath)
}

func (pnpFS *pnpFS) Chtimes(path string, mtime time.Time, atime time.Time) error {
	path, _, _ = resolveVirtual(path)

	fs, formattedPath, _ := getMatchingFS(pnpFS, path)
	return fs.Chtimes(formattedPath, mtime, atime)
}

func (pnpFS *pnpFS) Realpath(path string) string {
	path, hash, basePath := resolveVirtual(path)

	fs, formattedPath, zipPath := getMatchingFS(pnpFS, path)
	fullPath := zipPath + fs.Realpath(formattedPath)
	return makeVirtualPath(basePath, hash, fullPath)
}

func (pnpFS *pnpFS) Remove(path string) error {
	path, _, _ = resolveVirtual(path)

	fs, formattedPath, _ := getMatchingFS(pnpFS, path)
	return fs.Remove(formattedPath)
}

func (pnpFS *pnpFS) Stat(path string) vfs.FileInfo {
	path, _, _ = resolveVirtual(path)

	fs, formattedPath, _ := getMatchingFS(pnpFS, path)
	return fs.Stat(formattedPath)
}

func (pnpFS *pnpFS) UseCaseSensitiveFileNames() bool {
	// pnp fs is always case sensitive
	return true
}

func (pnpFS *pnpFS) WalkDir(root string, walkFn vfs.WalkDirFunc) error {
	root, hash, basePath := resolveVirtual(root)

	fs, formattedPath, zipPath := getMatchingFS(pnpFS, root)
	return fs.WalkDir(formattedPath, (func(path string, d vfs.DirEntry, err error) error {
		fullPath := zipPath + path
		return walkFn(makeVirtualPath(basePath, hash, fullPath), d, err)
	}))
}

func (pnpFS *pnpFS) WriteFile(path string, data string, writeByteOrderMark bool) error {
	path, _, _ = resolveVirtual(path)

	fs, formattedPath, zipPath := getMatchingFS(pnpFS, path)
	if zipPath != "" {
		panic("cannot write to zip file")
	}

	return fs.WriteFile(formattedPath, data, writeByteOrderMark)
}

func splitZipPath(path string) (string, string) {
	parts := strings.Split(path, ".zip/")
	if len(parts) < 2 {
		return path, "/"
	}
	return parts[0] + ".zip", "/" + parts[1]
}

func getMatchingFS(pnpFS *pnpFS, path string) (vfs.FS, string, string) {
	if !tspath.IsZipPath(path) {
		return pnpFS.fs, path, ""
	}

	zipPath, internalPath := splitZipPath(path)

	zipStat := pnpFS.fs.Stat(zipPath)
	if zipStat == nil {
		return pnpFS.fs, path, ""
	}

	var usedReader *zip.ReadCloser

	pnpFS.cacheReaderMutex.Lock()
	defer pnpFS.cacheReaderMutex.Unlock()

	cachedReader, ok := pnpFS.cachedZipReadersMap[zipPath]
	if ok {
		usedReader = cachedReader
	} else {
		zipReader, err := zip.OpenReader(zipPath)
		if err != nil {
			return pnpFS.fs, path, ""
		}

		usedReader = zipReader
		pnpFS.cachedZipReadersMap[zipPath] = usedReader
	}

	return iovfs.From(usedReader, pnpFS.fs.UseCaseSensitiveFileNames()), internalPath, zipPath
}

// Virtual paths are used to make different paths resolve to the same real file or folder, which is necessary in some cases when PnP is enabled
// See https://yarnpkg.com/advanced/lexicon#virtual-package and https://yarnpkg.com/advanced/pnpapi#resolvevirtual for more details
func resolveVirtual(path string) (realPath string, hash string, basePath string) {
	idx := strings.Index(path, "/__virtual__/")
	if idx == -1 {
		return path, "", ""
	}

	base := path[:idx]
	rest := path[idx+len("/__virtual__/"):]
	parts := strings.SplitN(rest, "/", 3)
	if len(parts) < 3 {
		// Not enough parts to match the pattern, return as is
		return path, "", ""
	}
	hash = parts[0]
	subpath := parts[2]
	depth, err := strconv.Atoi(parts[1])
	if err != nil || depth < 0 {
		// Invalid n, return as is
		return path, "", ""
	}

	basePath = path[:idx] + "/__virtual__"

	// Apply dirname n times to base
	for range depth {
		base = tspath.GetDirectoryPath(base)
	}
	// Join base and subpath
	if base == "/" {
		return "/" + subpath, hash, basePath
	}

	return tspath.CombinePaths(base, subpath), hash, basePath
}

func makeVirtualPath(basePath string, hash string, targetPath string) string {
	if basePath == "" || hash == "" {
		return targetPath
	}

	relativePath := tspath.GetRelativePathFromDirectory(
		tspath.GetDirectoryPath(basePath),
		targetPath,
		tspath.ComparePathsOptions{UseCaseSensitiveFileNames: true})

	segments := strings.Split(relativePath, "/")

	depth := 0
	for depth < len(segments) && segments[depth] == ".." {
		depth++
	}

	subPath := strings.Join(segments[depth:], "/")

	return path.Join(basePath, hash, strconv.Itoa(depth), subPath)
}

func (pnpFS *pnpFS) ClearCache() error {
	pnpFS.cacheReaderMutex.Lock()
	defer pnpFS.cacheReaderMutex.Unlock()

	for _, reader := range pnpFS.cachedZipReadersMap {
		err := reader.Close()
		if err != nil {
			return err
		}
	}

	pnpFS.cachedZipReadersMap = make(map[string]*zip.ReadCloser)

	return nil
}
