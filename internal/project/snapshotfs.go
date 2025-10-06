package project

import (
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/project/dirty"
	"github.com/microsoft/typescript-go/internal/sourcemap"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs"
	"github.com/microsoft/typescript-go/internal/vfs/cachedvfs"
	"github.com/zeebo/xxh3"
)

type FileSource interface {
	FS() vfs.FS
	GetFile(fileName string) FileHandle
}

var (
	_ FileSource = (*snapshotFSBuilder)(nil)
	_ FileSource = (*snapshotFS)(nil)
)

type snapshotFS struct {
	toPath    func(fileName string) tspath.Path
	fs        vfs.FS
	overlays  map[tspath.Path]*overlay
	diskFiles map[tspath.Path]*diskFile
}

func (s *snapshotFS) FS() vfs.FS {
	return s.fs
}

func (s *snapshotFS) GetFile(fileName string) FileHandle {
	if file, ok := s.overlays[s.toPath(fileName)]; ok {
		return file
	}
	if file, ok := s.diskFiles[s.toPath(fileName)]; ok {
		return file
	}
	return nil
}

func (s *snapshotFS) GetDocumentPositionMapper(fileName string) *sourcemap.DocumentPositionMapper {
	var sourceMapInfo *sourceMapInfo
	if file, ok := s.overlays[s.toPath(fileName)]; ok {
		sourceMapInfo = file.sourceMapInfo
	}
	if file, ok := s.diskFiles[s.toPath(fileName)]; ok {
		sourceMapInfo = file.sourceMapInfo
	}
	if sourceMapInfo != nil {
		if sourceMapInfo.documentMapper != nil {
			return sourceMapInfo.documentMapper.m
		}
		if sourceMapInfo.sourceMapPath != "" {
			return s.GetDocumentPositionMapper(sourceMapInfo.sourceMapPath)
		}
	}
	return nil
}

type snapshotFSBuilder struct {
	fs        vfs.FS
	overlays  *dirty.SyncMap[tspath.Path, *overlay]
	diskFiles *dirty.SyncMap[tspath.Path, *diskFile]
	toPath    func(string) tspath.Path
}

func newSnapshotFSBuilder(
	fs vfs.FS,
	overlays map[tspath.Path]*overlay,
	diskFiles map[tspath.Path]*diskFile,
	positionEncoding lsproto.PositionEncodingKind,
	toPath func(fileName string) tspath.Path,
) *snapshotFSBuilder {
	cachedFS := cachedvfs.From(fs)
	cachedFS.Enable()
	return &snapshotFSBuilder{
		fs:        cachedFS,
		overlays:  dirty.NewSyncMap(overlays, nil),
		diskFiles: dirty.NewSyncMap(diskFiles, nil),
		toPath:    toPath,
	}
}

func (s *snapshotFSBuilder) FS() vfs.FS {
	return s.fs
}

func (s *snapshotFSBuilder) Finalize() (*snapshotFS, bool) {
	diskFiles, diskChanged := s.diskFiles.Finalize()
	overlays, overlayChanged := s.overlays.Finalize()
	return &snapshotFS{
		fs:        s.fs,
		overlays:  overlays,
		diskFiles: diskFiles,
		toPath:    s.toPath,
	}, diskChanged || overlayChanged
}

type builderFileChanges struct {
	diskChanges    map[tspath.Path]*dirty.Change[*diskFile]
	overlayChanges map[tspath.Path]*dirty.Change[*overlay]
}

func (s *snapshotFSBuilder) Finalize2() (*snapshotFS, bool, *builderFileChanges) {
	diskFiles, diskChanged, diskChanges := s.diskFiles.Finalize2()
	overlays, overlayChanged, overlayChanges := s.overlays.Finalize2()
	return &snapshotFS{
		fs:        s.fs,
		overlays:  overlays,
		diskFiles: diskFiles,
		toPath:    s.toPath,
	}, diskChanged || overlayChanged, &builderFileChanges{diskChanges, overlayChanges}
}

func (s *snapshotFSBuilder) GetFile(fileName string) FileHandle {
	path := s.toPath(fileName)
	return s.GetFileByPath(fileName, path)
}

func (s *snapshotFSBuilder) GetFileByPath(fileName string, path tspath.Path) FileHandle {
	if entry, ok := s.overlays.Load(path); ok {
		return entry.Value()
	}
	entry, _ := s.diskFiles.LoadOrStore(path, &diskFile{fileBase: fileBase{fileName: fileName}, needsReload: true})
	if entry != nil {
		entry.Locked(func(entry dirty.Value[*diskFile]) {
			if entry.Value() != nil && !entry.Value().MatchesDiskText() {
				if content, ok := s.fs.ReadFile(fileName); ok {
					entry.Change(func(file *diskFile) {
						file.content = content
						file.hash = xxh3.Hash128([]byte(content))
						file.needsReload = false
						// This should not ever be needed while updating a snapshot with source maps
						file.sourceMapInfo = nil
						file.lineInfo = nil
					})
				} else {
					entry.Delete()
				}
			}
		})
	}
	if entry == nil || entry.Value() == nil {
		return nil
	}
	return entry.Value()
}

func (s *snapshotFSBuilder) markDirtyFiles(change FileChangeSummary) {
	for uri := range change.Changed.Keys() {
		path := s.toPath(uri.FileName())
		if entry, ok := s.diskFiles.Load(path); ok {
			entry.Change(func(file *diskFile) {
				file.needsReload = true
			})
		}
	}
	for uri := range change.Deleted.Keys() {
		path := s.toPath(uri.FileName())
		if entry, ok := s.diskFiles.Load(path); ok {
			entry.Change(func(file *diskFile) {
				file.needsReload = true
			})
		}
	}
}

// UseCaseSensitiveFileNames implements sourcemap.Host.
func (s *snapshotFSBuilder) UseCaseSensitiveFileNames() bool {
	return s.fs.UseCaseSensitiveFileNames()
}

// GetECMALineInfo implements sourcemap.Host.
func (s *snapshotFSBuilder) GetECMALineInfo(fileName string) *sourcemap.ECMALineInfo {
	if file := s.GetFile(fileName); file != nil {
		return file.ECMALineInfo()
	}
	return nil
}

// ReadFile implements sourcemap.Host.
func (s *snapshotFSBuilder) ReadFile(fileName string) (contents string, ok bool) {
	if file := s.GetFile(fileName); file != nil {
		return file.Content(), true
	}
	return "", false
}

// !!! TODO: consolidate this and `GetFileByPath`
func (s *snapshotFSBuilder) getDiskFileByPath(fileName string, path tspath.Path) *diskFile {
	entry, _ := s.diskFiles.LoadOrStore(path, &diskFile{fileBase: fileBase{fileName: fileName}, needsReload: true})
	if entry != nil {
		entry.Locked(func(entry dirty.Value[*diskFile]) {
			if entry.Value() != nil && !entry.Value().MatchesDiskText() {
				if content, ok := s.fs.ReadFile(fileName); ok {
					entry.Change(func(file *diskFile) {
						file.content = content
						file.hash = xxh3.Hash128([]byte(content))
						file.needsReload = false
						// This should not ever be needed while updating a snapshot with source maps
						file.sourceMapInfo = nil
						file.lineInfo = nil
					})
				} else {
					entry.Delete()
				}
			}
		})
	}
	if entry == nil || entry.Value() == nil {
		return nil
	}
	return entry.Value()
}

func (s *snapshotFSBuilder) computeDocumentPositionMapper(genFileName string) {
	genFilePath := s.toPath(genFileName)
	if entry, ok := s.overlays.Load(genFilePath); ok {
		file := entry.Value()
		if file == nil {
			return
		}
		// Source map information already computed
		if file.sourceMapInfo != nil {
			return
		}
		// Compute source map information
		url, isInline := sourcemap.GetSourceMapURL(s, genFileName)
		if isInline {
			// Store document mapper directly in disk file for an inline source map
			docMapper := sourcemap.ConvertDocumentToSourceMapper(s, url, genFileName)
			entry.Change(func(file *overlay) {
				file.sourceMapInfo = &sourceMapInfo{documentMapper: &documentMapper{m: docMapper}}
			})
		} else {
			// Store path to map file
			entry.Change(func(file *overlay) {
				file.sourceMapInfo = &sourceMapInfo{sourceMapPath: url}
			})
		}
		if url != "" {
			s.computeDocumentPositionMapperForMap(url)
		}
	}
	entry, ok := s.diskFiles.Load(genFilePath)
	if !ok {
		return
	}
	file := entry.Value()
	if file == nil {
		return
	}
	// Source map information already computed
	if file.sourceMapInfo != nil {
		return
	}
	// Compute source map information
	url, isInline := sourcemap.GetSourceMapURL(s, genFileName)
	if isInline {
		// Store document mapper directly in disk file for an inline source map
		docMapper := sourcemap.ConvertDocumentToSourceMapper(s, url, genFileName)
		entry.Change(func(file *diskFile) {
			file.sourceMapInfo = &sourceMapInfo{documentMapper: &documentMapper{m: docMapper}}
		})
	} else {
		// Store path to map file
		entry.Change(func(file *diskFile) {
			file.sourceMapInfo = &sourceMapInfo{sourceMapPath: url}
		})
	}
	if url != "" {
		s.computeDocumentPositionMapperForMap(url)
	}
}

func (s *snapshotFSBuilder) computeDocumentPositionMapperForMap(mapFileName string) {
	mapFilePath := s.toPath(mapFileName)
	mapFile := s.getDiskFileByPath(mapFileName, mapFilePath)
	if mapFile == nil {
		return
	}
	if mapFile.sourceMapInfo != nil {
		return
	}
	docMapper := sourcemap.ConvertDocumentToSourceMapper(s, mapFile.Content(), mapFileName)
	entry, _ := s.diskFiles.Load(mapFilePath)
	entry.Change(func(file *diskFile) {
		file.sourceMapInfo = &sourceMapInfo{documentMapper: &documentMapper{m: docMapper}}
	})
}

func (s *snapshotFSBuilder) applyDiskFileChanges(
	diskChanges map[tspath.Path]*dirty.Change[*diskFile],
) {
	for path, change := range diskChanges {
		if change.Deleted {
			if change.Old != nil {
				panic("Deleting files not supported")
			}
			// Failed file read
			continue
		}
		// New file
		if change.Old == nil {
			s.diskFiles.LoadOrStore(path, change.New)
			continue
		}
		// Updated file
		entry, _ := s.diskFiles.Load(path)
		if entry != nil {
			entry.Change(func(file *diskFile) {
				if file.Hash() == change.Old.Hash() {
					file.sourceMapInfo = change.New.sourceMapInfo
				}
			})
		}
	}
}
