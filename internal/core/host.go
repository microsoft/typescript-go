package core

import (
	"github.com/microsoft/typescript-go/internal/vfs"
)

type CompilerHost interface {
	FS() vfs.FS
	GetCurrentDirectory() string
	Trace(msg string)
}

type FileInfo struct {
	Name string
	Size int64
}

var _ CompilerHost = (*compilerHost)(nil)

type compilerHost struct {
	options          *CompilerOptions
	currentDirectory string
	fs               vfs.FS
}

func NewCompilerHost(options *CompilerOptions, currentDirectory string, fs vfs.FS) CompilerHost {
	h := &compilerHost{}
	h.options = options
	h.currentDirectory = currentDirectory
	h.fs = fs
	return h
}

func (h *compilerHost) FS() vfs.FS {
	return h.fs
}

func (h *compilerHost) GetCurrentDirectory() string {
	return h.currentDirectory
}

func (d *compilerHost) Trace(msg string) {
	//!!! TODO: implement
}
