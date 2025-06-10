package compiler

import (
	"strings"

	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/module"
	"github.com/microsoft/typescript-go/internal/modulespecifiers"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs"
)

type ProjectReferenceDtsFakingHost struct {
	projectReferenceFileMapper *ProjectReferenceFileMapper
	dtsDirectories             core.Set[tspath.Path]
	symlinkCache               modulespecifiers.SymlinkCache
}

var _ module.ResolutionHost = (*ProjectReferenceDtsFakingHost)(nil)

func (h *ProjectReferenceDtsFakingHost) FS() vfs.FS {
	return h
}

func (h *ProjectReferenceDtsFakingHost) GetCurrentDirectory() string {
	return h.projectReferenceFileMapper.opts.Host.GetCurrentDirectory()
}

func (h *ProjectReferenceDtsFakingHost) Trace(msg string) {
	h.projectReferenceFileMapper.opts.Host.Trace(msg)
}

// UseCaseSensitiveFileNames returns true if the file system is case-sensitive.
func (h *ProjectReferenceDtsFakingHost) UseCaseSensitiveFileNames() bool {
	return h.projectReferenceFileMapper.opts.Host.FS().UseCaseSensitiveFileNames()
}

// FileExists returns true if the file exists.
func (h *ProjectReferenceDtsFakingHost) FileExists(path string) bool {
	if h.projectReferenceFileMapper.opts.Host.FS().FileExists(path) {
		return true
	}
	if !tspath.IsDeclarationFileName(path) {
		return false
	}
	// Project references go to source file instead of .d.ts file
	return h.fileOrDirectoryExistsUsingSource(path /*isFile*/, true)
}

func (h *ProjectReferenceDtsFakingHost) ReadFile(path string) (contents string, ok bool) {
	// Dont need to override as we cannot mimick read file
	return h.projectReferenceFileMapper.opts.Host.FS().ReadFile(path)
}

func (h *ProjectReferenceDtsFakingHost) WriteFile(path string, data string, writeByteOrderMark bool) error {
	panic("should not be called by resolver")
}

// Removes `path` and all its contents. Will return the first error it encounters.
func (h *ProjectReferenceDtsFakingHost) Remove(path string) error {
	panic("should not be called by resolver")
}

// DirectoryExists returns true if the path is a directory.
func (h *ProjectReferenceDtsFakingHost) DirectoryExists(path string) bool {
	if h.projectReferenceFileMapper.opts.Host.FS().DirectoryExists(path) {
		h.handleDirectoryCouldBeSymlink(path)
		return true
	}
	return h.fileOrDirectoryExistsUsingSource(path /*isFile*/, false)
}

// GetAccessibleEntries returns the files/directories in the specified directory.
// If any entry is a symlink, it will be followed.
func (h *ProjectReferenceDtsFakingHost) GetAccessibleEntries(path string) vfs.Entries {
	panic("should not be called by resolver")
}

func (h *ProjectReferenceDtsFakingHost) Stat(path string) vfs.FileInfo {
	panic("should not be called by resolver")
}

// WalkDir walks the file tree rooted at root, calling walkFn for each file or directory in the tree.
// It is has the same behavior as [fs.WalkDir], but with paths as [string].
func (h *ProjectReferenceDtsFakingHost) WalkDir(root string, walkFn vfs.WalkDirFunc) error {
	panic("should not be called by resolver")
}

// Realpath returns the "real path" of the specified path,
// following symlinks and correcting filename casing.
func (h *ProjectReferenceDtsFakingHost) Realpath(path string) string {
	result, ok := h.symlinkCache.SymlinkedFiles().Load(h.toPath(path))
	if ok {
		return result
	}
	return h.projectReferenceFileMapper.opts.Host.FS().Realpath(path)
}

func (h *ProjectReferenceDtsFakingHost) toPath(path string) tspath.Path {
	return tspath.ToPath(path, h.GetCurrentDirectory(), h.UseCaseSensitiveFileNames())
}

func (h *ProjectReferenceDtsFakingHost) handleDirectoryCouldBeSymlink(directory string) {
	if tspath.ContainsIgnoredPath(directory) {
		return
	}

	// Because we already watch node_modules, handle symlinks in there
	if !strings.Contains(directory, "/node_modules/") {
		return
	}
	directoryPath := tspath.Path(tspath.EnsureTrailingDirectorySeparator(string(h.toPath(directory))))
	if _, ok := h.symlinkCache.SymlinkedDirectories().Load(directoryPath); ok {
		return
	}

	realDirectory := h.projectReferenceFileMapper.opts.Host.FS().Realpath(directory)
	var realPath tspath.Path
	var symlinkedDirectory *modulespecifiers.SymlinkedDirectory
	if realDirectory == directory {
		// not symlinked
		symlinkedDirectory = nil
	} else if realPath = tspath.Path(tspath.EnsureTrailingDirectorySeparator(string(h.toPath(realDirectory)))); realPath == directoryPath {
		// not symlinked
		symlinkedDirectory = nil
	} else {
		symlinkedDirectory = &modulespecifiers.SymlinkedDirectory{
			Real:     tspath.EnsureTrailingDirectorySeparator(realDirectory),
			RealPath: realPath,
		}
	}
	h.symlinkCache.SetSymlinkedDirectory(directory, directoryPath, symlinkedDirectory)
}

func (h *ProjectReferenceDtsFakingHost) fileOrDirectoryExistsUsingSource(fileOrDirectory string, isFile bool) bool {
	fileOrDirectoryExistsUsingSource := core.IfElse(isFile, h.fileExistsIfProjectReferenceDts, h.directoryExistsIfProjectReferenceDeclDir)
	// Check current directory or file
	result := fileOrDirectoryExistsUsingSource(fileOrDirectory)
	if result != core.TSUnknown {
		return result == core.TSTrue
	}

	symlinkedDirectories := h.symlinkCache.SymlinkedDirectories()
	if symlinkedDirectories.Size() == 0 {
		return false
	}
	fileOrDirectoryPath := h.toPath(fileOrDirectory)
	if !strings.Contains(string(fileOrDirectoryPath), "/node_modules/") {
		return false
	}
	if isFile {
		_, ok := h.symlinkCache.SymlinkedFiles().Load(fileOrDirectoryPath)
		if ok {
			return true
		}
	}

	// If it contains node_modules check if its one of the symlinked path we know of
	var exists bool
	symlinkedDirectories.Range(func(directoryPath tspath.Path, symlinkedDirectory *modulespecifiers.SymlinkedDirectory) bool {
		if symlinkedDirectory == nil {
			return true
		}

		relative, hasPrefix := strings.CutPrefix(string(fileOrDirectoryPath), string(directoryPath))
		if !hasPrefix {
			return true
		}
		if exists = fileOrDirectoryExistsUsingSource(string(symlinkedDirectory.RealPath) + relative).IsTrue(); exists {
			if isFile {
				// Store the real path for the file
				absolutePath := tspath.GetNormalizedAbsolutePath(fileOrDirectory, h.GetCurrentDirectory())
				h.symlinkCache.SetSymlinkedFile(
					fileOrDirectoryPath,
					symlinkedDirectory.Real+absolutePath[len(directoryPath):],
				)
			}
			return false
		}
		return true
	})
	return exists
}

func (h *ProjectReferenceDtsFakingHost) fileExistsIfProjectReferenceDts(file string) core.Tristate {
	source := h.projectReferenceFileMapper.getSourceAndProjectReference(h.toPath(file))
	if source != nil {
		return core.IfElse(h.projectReferenceFileMapper.opts.Host.FS().FileExists(source.Source), core.TSTrue, core.TSFalse)
	}
	return core.TSUnknown
}

func (h *ProjectReferenceDtsFakingHost) directoryExistsIfProjectReferenceDeclDir(dir string) core.Tristate {
	dirPath := h.toPath(dir)
	dirPathWithTrailingDirectorySeparator := dirPath + "/"
	for declDirPath := range h.dtsDirectories.Keys() {
		if dirPath == declDirPath ||
			// Any parent directory of declaration dir
			strings.HasPrefix(string(declDirPath), string(dirPathWithTrailingDirectorySeparator)) ||
			// Any directory inside declaration dir
			strings.HasPrefix(string(dirPath), string(declDirPath)+"/") {
			return core.TSTrue
		}
	}
	return core.TSUnknown
}
