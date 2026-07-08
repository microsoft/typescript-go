// Package zipoverlay wraps a vfs.FS so the compiler can read files that live
// inside Yarn PnP cache archives. Under PnP a package's on-disk location can be a
// path whose middle component is a .zip file, e.g.
//
//	/proj/.yarn/cache/react-npm-18.3.1-<hash>.zip/node_modules/react/index.d.ts
//
// The OS reports the .zip as a regular file, so every path inside it fails
// FileExists/DirectoryExists and the module resolver's parent-directory checks
// abort. This overlay intercepts any path that resolves into a cache archive and
// serves it from the archive, reporting the archive and its interior directories
// as directories so those checks pass.
//
// It also dereferences Yarn's __virtual__/<hash>/<n> paths (peer-dependency
// virtualization) to their backing location before touching the filesystem — a
// cache .zip for an npm package, or a plain on-disk directory for a
// workspace/unplugged package. The virtual path itself is preserved everywhere
// above the filesystem (the resolver's locator table is keyed by it), so only the
// actual read is redirected. All other paths delegate to the inner FS unchanged.
//
// Zip-internal lookups are case-sensitive. On a case-insensitive host a module
// specifier whose case differs from the archived entry will not resolve; that is
// a known limitation (see #460).
package zipoverlay

import (
	"archive/zip"
	"io"
	"io/fs"
	"strings"
	"sync"
	"time"

	"github.com/microsoft/typescript-go/internal/pnp"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs"
)

const zipMarker = ".zip/"

type overlay struct {
	inner vfs.FS

	mu      sync.Mutex
	archive map[string]*archiveIndex // keyed by on-disk .zip path
}

var _ vfs.FS = (*overlay)(nil)

// Wrap returns a vfs.FS that serves paths inside Yarn cache .zip archives from
// those archives (dereferencing PnP virtual paths first) and delegates everything
// else to inner.
func Wrap(inner vfs.FS) vfs.FS {
	return &overlay{
		inner:   inner,
		archive: map[string]*archiveIndex{},
	}
}

// split resolves a path for the overlay. It dereferences a PnP virtual segment,
// then reports whether the result addresses a real .zip archive. target is the
// path to delegate to the inner FS when ok is false (the dereferenced path for a
// virtual location backed by a plain directory, or the original path otherwise).
func (o *overlay) split(path string) (zipPath string, internal string, target string, ok bool) {
	target = path
	if strings.Contains(path, "__virtual__") || strings.Contains(path, "$$virtual") {
		target = pnp.DerefVirtualPath(path)
	} else if !strings.Contains(path, ".zip") {
		// Fast path: not virtual and no archive component. The overlay wraps the OS
		// FS unconditionally, so this runs on every filesystem probe in every
		// project; keep it a single scan.
		return "", "", path, false
	}
	if idx := strings.Index(target, zipMarker); idx >= 0 {
		zipPath = target[:idx+len(zipMarker)-1] // include ".zip", drop the "/"
		internal = strings.Trim(target[idx+len(zipMarker):], "/")
	} else if strings.HasSuffix(target, ".zip") {
		zipPath = target
		internal = ""
	} else {
		// A virtual path backed by a plain directory (workspace / unplugged
		// package): delegate the dereferenced path to the inner FS.
		return "", "", target, false
	}
	if !o.inner.FileExists(zipPath) {
		return "", "", target, false
	}
	return zipPath, internal, target, true
}

// index lazily opens and indexes the archive at zipPath.
func (o *overlay) index(zipPath string) *archiveIndex {
	o.mu.Lock()
	defer o.mu.Unlock()
	if a, ok := o.archive[zipPath]; ok {
		return a
	}
	a := o.load(zipPath)
	// Only cache a successful load. A nil load means the .zip exists but could not
	// be read yet (e.g. a still-writing install); caching nil would poison the
	// path for the life of a long-running LSP process even after the archive
	// becomes valid.
	if a != nil {
		o.archive[zipPath] = a
	}
	return a
}

func (o *overlay) load(zipPath string) *archiveIndex {
	contents, ok := o.inner.ReadFile(zipPath)
	if !ok {
		return nil
	}
	r, err := zip.NewReader(strings.NewReader(contents), int64(len(contents)))
	if err != nil {
		return nil
	}
	a := &archiveIndex{
		files: map[string]*zip.File{},
		dirs:  map[string]bool{"": true},
	}
	for _, f := range r.File {
		name := strings.Trim(f.Name, "/")
		if name == "" {
			continue
		}
		if strings.HasSuffix(f.Name, "/") {
			a.dirs[name] = true
		} else {
			a.files[name] = f
		}
		// Register every ancestor directory of this entry.
		for i := range len(name) {
			if name[i] == '/' {
				a.dirs[name[:i]] = true
			}
		}
	}
	return a
}

type archiveIndex struct {
	files map[string]*zip.File // interior path -> file
	dirs  map[string]bool      // interior dir paths (plus "")
}

func (o *overlay) UseCaseSensitiveFileNames() bool {
	return o.inner.UseCaseSensitiveFileNames()
}

func (o *overlay) FileExists(path string) bool {
	zipPath, internal, target, ok := o.split(path)
	if !ok {
		return o.inner.FileExists(target)
	}
	if a := o.index(zipPath); a != nil {
		_, isFile := a.files[internal]
		return isFile
	}
	return false
}

func (o *overlay) DirectoryExists(path string) bool {
	zipPath, internal, target, ok := o.split(path)
	if !ok {
		return o.inner.DirectoryExists(target)
	}
	if a := o.index(zipPath); a != nil {
		return a.dirs[internal]
	}
	return false
}

func (o *overlay) ReadFile(path string) (string, bool) {
	zipPath, internal, target, ok := o.split(path)
	if !ok {
		return o.inner.ReadFile(target)
	}
	a := o.index(zipPath)
	if a == nil {
		return "", false
	}
	f, isFile := a.files[internal]
	if !isFile {
		return "", false
	}
	rc, err := f.Open()
	if err != nil {
		return "", false
	}
	defer rc.Close()
	b, err := io.ReadAll(rc)
	if err != nil {
		return "", false
	}
	return string(b), true
}

func (o *overlay) Stat(path string) vfs.FileInfo {
	zipPath, internal, target, ok := o.split(path)
	if !ok {
		return o.inner.Stat(target)
	}
	a := o.index(zipPath)
	if a == nil {
		return nil
	}
	if a.dirs[internal] {
		return &fileInfo{name: tspath.GetBaseFileName(internal), mode: fs.ModeDir}
	}
	if f, isFile := a.files[internal]; isFile {
		return &fileInfo{name: tspath.GetBaseFileName(internal), size: int64(f.UncompressedSize64)}
	}
	return nil
}

func (o *overlay) GetAccessibleEntries(path string) vfs.Entries {
	zipPath, internal, target, ok := o.split(path)
	if !ok {
		return o.inner.GetAccessibleEntries(target)
	}
	var entries vfs.Entries
	a := o.index(zipPath)
	if a == nil {
		return entries
	}
	prefix := internal
	if prefix != "" {
		prefix += "/"
	}
	for name := range a.files {
		if child, isChild := immediateChild(prefix, name); isChild {
			entries.Files = append(entries.Files, child)
		}
	}
	seenDir := map[string]bool{}
	for name := range a.dirs {
		if name == "" {
			continue
		}
		if child, isChild := immediateChild(prefix, name); isChild && !seenDir[child] {
			seenDir[child] = true
			entries.Directories = append(entries.Directories, child)
		}
	}
	return entries
}

func (o *overlay) WalkDir(root string, walkFn vfs.WalkDirFunc) error {
	zipPath, internal, target, ok := o.split(root)
	if !ok {
		return o.inner.WalkDir(target, walkFn)
	}
	a := o.index(zipPath)
	if a == nil {
		return nil
	}
	return o.walk(zipPath, a, internal, walkFn)
}

func (o *overlay) walk(zipPath string, a *archiveIndex, dir string, walkFn vfs.WalkDirFunc) error {
	full := zipPath + "/" + dir
	if dir == "" {
		full = zipPath
	}
	if err := walkFn(full, fs.FileInfoToDirEntry(&fileInfo{name: tspath.GetBaseFileName(dir), mode: fs.ModeDir}), nil); err != nil {
		if err == fs.SkipDir || err == fs.SkipAll { //nolint:errorlint
			return nil
		}
		return err
	}
	entries := o.GetAccessibleEntries(full)
	for _, f := range entries.Files {
		child := joinInternal(dir, f)
		size := int64(0)
		if zf, ok := a.files[child]; ok {
			size = int64(zf.UncompressedSize64)
		}
		if err := walkFn(zipPath+"/"+child, fs.FileInfoToDirEntry(&fileInfo{name: f, size: size}), nil); err != nil {
			if err == fs.SkipAll { //nolint:errorlint
				return nil
			}
			if err == fs.SkipDir { //nolint:errorlint
				continue
			}
			return err
		}
	}
	for _, d := range entries.Directories {
		if err := o.walk(zipPath, a, joinInternal(dir, d), walkFn); err != nil {
			return err
		}
	}
	return nil
}

func (o *overlay) Realpath(path string) string {
	// Zip-internal and PnP virtual paths are already canonical for the compiler's
	// purposes (cache filenames are content-addressed; there are no symlinks inside
	// an archive) and, crucially, must be returned unchanged: the resolver's PnP
	// locator table is keyed by the virtual path, so canonicalizing to the real
	// backing location here would desync findLocator for imports made from inside a
	// virtualized package.
	if strings.Contains(path, zipMarker) ||
		strings.HasSuffix(path, ".zip") ||
		strings.Contains(path, "__virtual__") ||
		strings.Contains(path, "$$virtual") {
		return path
	}
	return o.inner.Realpath(path)
}

// Mutating operations never target zip-internal or virtual paths in this
// compiler; delegate the original path unchanged.
func (o *overlay) WriteFile(path string, data string) error { return o.inner.WriteFile(path, data) }

func (o *overlay) AppendFile(path string, data string) error {
	return o.inner.AppendFile(path, data)
}
func (o *overlay) Remove(path string) error { return o.inner.Remove(path) }
func (o *overlay) Chtimes(path string, aTime time.Time, mTime time.Time) error {
	return o.inner.Chtimes(path, aTime, mTime)
}

// immediateChild reports whether name is a direct child of prefix ("" = root),
// returning the child's base name.
func immediateChild(prefix, name string) (string, bool) {
	if !strings.HasPrefix(name, prefix) {
		return "", false
	}
	rest := name[len(prefix):]
	if rest == "" || strings.Contains(rest, "/") {
		return "", false
	}
	return rest, true
}

func joinInternal(dir, child string) string {
	if dir == "" {
		return child
	}
	return dir + "/" + child
}

type fileInfo struct {
	name string
	size int64
	mode fs.FileMode
}

var (
	_ fs.FileInfo = (*fileInfo)(nil)
	_ fs.DirEntry = (*fileInfo)(nil)
)

func (fi *fileInfo) IsDir() bool                { return fi.mode.IsDir() }
func (fi *fileInfo) ModTime() time.Time         { return time.Time{} }
func (fi *fileInfo) Mode() fs.FileMode          { return fi.mode }
func (fi *fileInfo) Name() string               { return fi.name }
func (fi *fileInfo) Size() int64                { return fi.size }
func (fi *fileInfo) Sys() any                   { return nil }
func (fi *fileInfo) Info() (fs.FileInfo, error) { return fi, nil }
func (fi *fileInfo) Type() fs.FileMode          { return fi.mode.Type() }
