package project

import "github.com/microsoft/typescript-go/internal/vfs"

type ProjecServicetHost interface {
	FS() vfs.FS
	GetCurrentDirectory() string
	NewLine() string
	Trace(msg string)
}
