package vfstest

import (
	"fmt"
	"io/fs"
	"path"
	"strings"
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

// WithSensitivity converts a [fstest.MapFS] to one with the specified case sensitivity.
// The paths given in the map are treated as the "real paths".
//
// If useCaseSensitiveFileNames is true, the map is returned as-is.
func WithSensitivity(m fstest.MapFS, useCaseSensitiveFileNames bool) fs.FS {
	if useCaseSensitiveFileNames {
		return m
	}

	// Create all missing intermediate directories so we can attach the real path to each of them.
	newFiles := make(fstest.MapFS)
	for p := range m {
		curr := ""
		remaining := p

		for remaining != "" {
			before, after, _ := strings.Cut(remaining, "/")
			if curr == "" {
				curr = before
			} else {
				curr = curr + "/" + before
			}
			remaining = after

			if _, ok := m[curr]; !ok {
				newFiles[curr] = &fstest.MapFile{
					Mode: fs.ModeDir | 0555,
				}
			}
		}
	}

	newM := make(fstest.MapFS, len(m)+len(newFiles))
	for k, v := range m {
		newM[k] = v
	}
	for k, v := range newFiles {
		newM[k] = v
	}

	mp := make(fstest.MapFS, len(newM))
	for path, file := range newM {
		canonical := tspath.GetCanonicalFileName(path, useCaseSensitiveFileNames)
		if other, ok := mp[canonical]; ok {
			otherPath := other.Sys.(*sys).realPath
			panic(fmt.Sprintf("duplicate path: %q and %q have the same canonical path", path, otherPath))
		}
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
		info := must(entry.Info())
		newInfo, ok := convertInfo(info)
		if !ok {
			panic(fmt.Sprintf("unexpected synthesized dir: %q", info.Name()))
		}

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

	info := must(f.Stat())

	newInfo, ok := convertInfo(info)
	if !ok {
		// This is a synthesized dir.
		if name != "." {
			panic(fmt.Sprintf("unexpected synthesized dir: %q", name))
		}

		f := f.(fs.ReadDirFile)

		return &readDirFile{
			ReadDirFile: f,
			fileInfo: &fileInfo{
				FileInfo: info,
				sys:      info.Sys(),
				realPath: ".",
			},
		}, nil
	}

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

func convertInfo(info fs.FileInfo) (*fileInfo, bool) {
	sys, ok := info.Sys().(*sys)
	if !ok {
		return nil, false
	}
	return &fileInfo{
		FileInfo: info,
		sys:      sys.original,
		realPath: sys.realPath,
	}, true
}

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
