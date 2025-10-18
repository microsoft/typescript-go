package zipvfs

import (
	"archive/zip"
	"io"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/microsoft/typescript-go/internal/vfs"
)

type EntryKind uint8

const (
	DirEntry  EntryKind = 1
	FileEntry EntryKind = 2
)

type ZipFS struct {
	inner vfs.FS

	zipFilesMutex sync.Mutex
	zipFiles      map[string]*zipFile
}

type zipFile struct {
	reader *zip.ReadCloser
	err    error

	dirs  map[string]*compressedDir
	files map[string]*compressedFile
	once  sync.Once
}

type compressedDir struct {
	entries map[string]EntryKind
	path    string
	mutex   sync.Mutex
}

type compressedFile struct {
	compressed *zip.File

	// The file is decompressed lazily
	mutex    sync.Mutex
	contents string
	err      error
	wasRead  bool
}

func From(baseFS vfs.FS) *ZipFS {
	return &ZipFS{
		inner:    baseFS,
		zipFiles: make(map[string]*zipFile),
	}
}

func (fs *ZipFS) checkForZip(path string, kind EntryKind) (*zipFile, string) {
	var zipPath string
	var pathTail string

	if before, after, ok := strings.Cut(path, ".zip/"); ok {
		zipPath = before + ".zip"
		pathTail = after
	} else if kind == DirEntry && strings.HasSuffix(path, ".zip") {
		zipPath = path
	} else {
		return nil, ""
	}

	// If there is one, then check whether it's a file on the file system or not
	fs.zipFilesMutex.Lock()
	archive := fs.zipFiles[zipPath]
	if archive != nil {
		fs.zipFilesMutex.Unlock()
		archive.once.Do(func() {
			// wait if another goroutine is initializing the archive
		})
	} else {
		archive = &zipFile{}
		archive.once.Do(func() {
			fs.zipFiles[zipPath] = archive
			fs.zipFilesMutex.Unlock()

			tryToReadZipArchive(zipPath, archive)
		})
	}

	if archive.err != nil {
		return nil, ""
	}
	return archive, pathTail
}

func tryToReadZipArchive(zipPath string, archive *zipFile) {
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		archive.err = err
		return
	}

	dirs := make(map[string]*compressedDir)
	files := make(map[string]*compressedFile)
	seeds := []string{}

	// Build an index of all files in the archive
	for _, file := range reader.File {
		baseName := strings.TrimSuffix(file.Name, "/")
		dirPath := ""
		if slash := strings.LastIndexByte(baseName, '/'); slash != -1 {
			dirPath = baseName[:slash]
			baseName = baseName[slash+1:]
		}
		if file.FileInfo().IsDir() {
			// Handle a directory
			lowerDir := strings.ToLower(dirPath)
			if _, ok := dirs[lowerDir]; !ok {
				dir := &compressedDir{
					path:    dirPath,
					entries: make(map[string]EntryKind),
				}

				// List the same directory both with and without the slash
				dirs[lowerDir] = dir
				dirs[lowerDir+"/"] = dir
				seeds = append(seeds, lowerDir)
			}
		} else {
			// Handle a file
			files[strings.ToLower(file.Name)] = &compressedFile{compressed: file}
			lowerDir := strings.ToLower(dirPath)
			dir, ok := dirs[lowerDir]
			if !ok {
				dir = &compressedDir{
					path:    dirPath,
					entries: make(map[string]EntryKind),
				}

				// List the same directory both with and without the slash
				dirs[lowerDir] = dir
				dirs[lowerDir+"/"] = dir
				seeds = append(seeds, lowerDir)
			}
			dir.entries[baseName] = FileEntry
		}
	}

	// Populate child directories
	for _, baseName := range seeds {
		for baseName != "" {
			dirPath := ""
			if slash := strings.LastIndexByte(baseName, '/'); slash != -1 {
				dirPath = baseName[:slash]
				baseName = baseName[slash+1:]
			}
			lowerDir := strings.ToLower(dirPath)
			dir, ok := dirs[lowerDir]
			if !ok {
				dir = &compressedDir{
					path:    dirPath,
					entries: make(map[string]EntryKind),
				}

				// List the same directory both with and without the slash
				dirs[lowerDir] = dir
				dirs[lowerDir+"/"] = dir
			}
			dir.entries[baseName] = DirEntry
			baseName = dirPath
		}
	}

	archive.dirs = dirs
	archive.files = files
	archive.reader = reader
}

func (fs *ZipFS) UseCaseSensitiveFileNames() bool {
	return fs.inner.UseCaseSensitiveFileNames()
}

func (fs *ZipFS) FileExists(path string) bool {
	path = mangleYarnPnPVirtualPath(path)

	if fs.inner.FileExists(path) {
		return true
	}

	zip, pathTail := fs.checkForZip(path, FileEntry)
	if zip == nil {
		return false
	}

	_, ok := zip.files[strings.ToLower(pathTail)]
	return ok
}

func (fs *ZipFS) ReadFile(path string) (contents string, ok bool) {
	path = mangleYarnPnPVirtualPath(path)

	contents, ok = fs.inner.ReadFile(path)
	if ok {
		return contents, ok
	}

	// If the file doesn't exist, try reading from an enclosing zip archive
	zip, pathTail := fs.checkForZip(path, FileEntry)
	if zip == nil {
		return "", false
	}

	// Does the zip archive have this file?
	file, ok := zip.files[strings.ToLower(pathTail)]
	if !ok {
		return "", false
	}

	// Check whether it has already been read
	file.mutex.Lock()
	defer file.mutex.Unlock()

	if file.wasRead {
		return file.contents, file.err == nil
	}
	file.wasRead = true

	// If not, try to open it
	reader, err := file.compressed.Open()
	if err != nil {
		file.err = err
		return "", err == nil
	}
	defer reader.Close()

	// Then try to read it
	bytes, err := io.ReadAll(reader)
	if err != nil {
		file.err = err
		return "", false
	}

	file.contents = string(bytes)

	return file.contents, true
}

func (fs *ZipFS) WriteFile(path string, data string, writeByteOrderMark bool) error {
	err := fs.inner.WriteFile(path, data, writeByteOrderMark)
	if err != nil {
		fs.tryZipFileAssertion(path)
		return err
	}
	return nil
}

func (fs *ZipFS) Remove(path string) error {
	err := fs.inner.Remove(path)
	if err != nil {
		fs.tryZipFileAssertion(path)
		return err
	}
	return nil
}

func (fs *ZipFS) Chtimes(path string, aTime time.Time, mTime time.Time) error {
	err := fs.inner.Chtimes(path, aTime, mTime)
	if err != nil {
		fs.tryZipFileAssertion(path)
		return err
	}
	return nil
}

func (fs *ZipFS) DirectoryExists(path string) bool {
	path = mangleYarnPnPVirtualPath(path)

	if fs.inner.DirectoryExists(path) {
		return true
	}

	zip, pathTail := fs.checkForZip(path, DirEntry)
	if zip == nil {
		return false
	}

	_, ok := zip.dirs[strings.ToLower(pathTail)]
	return ok
}

func (fs *ZipFS) GetAccessibleEntries(path string) vfs.Entries {
	path = mangleYarnPnPVirtualPath(path)

	entries := fs.inner.GetAccessibleEntries(path)
	if len(entries.Files) > 0 || len(entries.Directories) > 0 {
		return entries
	}

	zip, pathTail := fs.checkForZip(path, DirEntry)
	if zip == nil {
		return entries
	}

	dir, ok := zip.dirs[strings.ToLower(pathTail)]
	if !ok {
		return entries
	}

	var files []string
	var dirs []string

	dir.mutex.Lock()
	defer dir.mutex.Unlock()

	for name, kind := range dir.entries {
		switch kind {
		case FileEntry:
			files = append(files, name)
		case DirEntry:
			dirs = append(dirs, name)
		}
	}

	return vfs.Entries{
		Files:       files,
		Directories: dirs,
	}
}

func (fs *ZipFS) Stat(path string) vfs.FileInfo {
	return fs.inner.Stat(path)
}

func (fs *ZipFS) WalkDir(root string, walkFn vfs.WalkDirFunc) error {
	return fs.inner.WalkDir(root, walkFn)
}

func (fs *ZipFS) Realpath(path string) string {
	return fs.inner.Realpath(path)
}

func (fs *ZipFS) tryZipFileAssertion(path string) {
	zip, _ := fs.checkForZip(path, FileEntry)
	if zip == nil {
		return
	}
	panic("do not use this method for zip file: " + path)
}

func ParseYarnPnPVirtualPath(path string) (string, string, bool) {
	i := 0

	for {
		start := i
		slash := strings.IndexAny(path[i:], "/\\")
		if slash == -1 {
			break
		}
		i += slash + 1

		// Replace the segments "__virtual__/<segment>/<n>" with N times the ".."
		// operation. Note: The "__virtual__" folder name appeared with Yarn 3.0.
		// Earlier releases used "$$virtual", but it was changed after discovering
		// that this pattern triggered bugs in software where paths were used as
		// either regexps or replacement. For example, "$$" found in the second
		// parameter of "String.prototype.replace" silently turned into "$".
		if segment := path[start : i-1]; segment == "__virtual__" || segment == "$$virtual" {
			if slash := strings.IndexAny(path[i:], "/\\"); slash != -1 {
				var count string
				var suffix string
				j := i + slash + 1

				// Find the range of the count
				if slash := strings.IndexAny(path[j:], "/\\"); slash != -1 {
					count = path[j : j+slash]
					suffix = path[j+slash:]
				} else {
					count = path[j:]
				}

				// Parse the count
				if n, err := strconv.ParseInt(count, 10, 64); err == nil {
					prefix := path[:start]

					// Apply N times the ".." operator
					for n > 0 && (strings.HasSuffix(prefix, "/") || strings.HasSuffix(prefix, "\\")) {
						slash := strings.LastIndexAny(prefix[:len(prefix)-1], "/\\")
						if slash == -1 {
							break
						}
						prefix = prefix[:slash+1]
						n--
					}

					// Make sure the prefix and suffix work well when joined together
					if suffix == "" && strings.IndexAny(prefix, "/\\") != strings.LastIndexAny(prefix, "/\\") {
						prefix = prefix[:len(prefix)-1]
					} else if prefix == "" {
						prefix = "."
					} else if strings.HasPrefix(suffix, "/") || strings.HasPrefix(suffix, "\\") {
						suffix = suffix[1:]
					}

					return prefix, suffix, true
				}
			}
		}
	}

	return "", "", false
}

func mangleYarnPnPVirtualPath(path string) string {
	if prefix, suffix, ok := ParseYarnPnPVirtualPath(path); ok {
		return prefix + suffix
	}
	return path
}
