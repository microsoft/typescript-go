package vfstest

import (
	"io/fs"
	"path"
	"testing/fstest"

	"github.com/microsoft/typescript-go/internal/tspath"
)

type mapFS struct {
	m                         fstest.MapFS
	useCaseSensitiveFileNames bool
}

type sys struct {
	original any
	realPath string
}

func ToMapFS(m fstest.MapFS, useCaseSensitiveFileNames bool) fs.FS {
	mp := make(fstest.MapFS, len(m))
	for path, file := range m {
		canonical := tspath.GetCanonicalFileName(path, useCaseSensitiveFileNames)
		fileCopy := *file
		fileCopy.Sys = &sys{
			original: fileCopy.Sys,
			realPath: path,
		}
		mp[canonical] = &fileCopy
	}
	return &mapFS{
		m:                         mp,
		useCaseSensitiveFileNames: useCaseSensitiveFileNames,
	}
}

type file struct {
	fs.File
	fileInfo *fileInfo
}

func (f *file) Stat() (fs.FileInfo, error) {
	return f.fileInfo, nil
}

type readDirFile struct {
	fs.ReadDirFile
	fileInfo *fileInfo
}

func (f *readDirFile) Stat() (fs.FileInfo, error) {
	return f.fileInfo, nil
}

func (f *readDirFile) ReadDir(n int) ([]fs.DirEntry, error) {
	list, err := f.ReadDirFile.ReadDir(n)
	if err != nil {
		return nil, err
	}

	entries := make([]fs.DirEntry, len(list))
	for i, entry := range list {
		info, err := entry.Info()
		if err != nil {
			return nil, err
		}
		newInfo := convertInfo(info)
		entries[i] = &dirEntry{
			DirEntry: entry,
			fileInfo: newInfo,
		}
	}

	return entries, nil
}

type dirEntry struct {
	fs.DirEntry
	fileInfo *fileInfo
}

func (e *dirEntry) Name() string {
	return path.Base(e.fileInfo.realPath)
}

func (e *dirEntry) Info() (fs.FileInfo, error) {
	return e.fileInfo, nil
}

type fileInfo struct {
	fs.FileInfo
	sys      any
	realPath string
}

func (fi *fileInfo) Name() string {
	return path.Base(fi.realPath)
}

func (fi *fileInfo) Sys() any {
	return fi.sys
}

func (m *mapFS) Open(name string) (fs.File, error) {
	f, err := m.m.Open(tspath.GetCanonicalFileName(name, m.useCaseSensitiveFileNames))
	if err != nil {
		return nil, err
	}

	info, err := f.Stat()
	if err != nil {
		return nil, err
	}

	newInfo := convertInfo(info)

	if f, ok := f.(fs.ReadDirFile); ok {
		return &readDirFile{
			ReadDirFile: f,
			fileInfo:    newInfo,
		}, nil
	}

	return &file{
		File:     f,
		fileInfo: newInfo,
	}, nil
}

func convertInfo(info fs.FileInfo) *fileInfo {
	sys := info.Sys().(*sys)
	return &fileInfo{
		FileInfo: info,
		sys:      sys.original,
		realPath: sys.realPath,
	}
}
