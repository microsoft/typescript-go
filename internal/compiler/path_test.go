package compiler

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestNormalizeSlashes(t *testing.T) {
	assert.Equal(t, normalizeSlashes("a"), "a")
	assert.Equal(t, normalizeSlashes("a/b"), "a/b")
	assert.Equal(t, normalizeSlashes("a\\b"), "a/b")
	assert.Equal(t, normalizeSlashes("\\\\server\\path"), "//server/path")
}

func TestGetRootLength(t *testing.T) {
	assert.Equal(t, getRootLength("a"), 0)
	assert.Equal(t, getRootLength("/"), 1)
	assert.Equal(t, getRootLength("/path"), 1)
	assert.Equal(t, getRootLength("c:"), 2)
	assert.Equal(t, getRootLength("c:d"), 0)
	assert.Equal(t, getRootLength("c:/"), 3)
	assert.Equal(t, getRootLength("c:\\"), 3)
	assert.Equal(t, getRootLength("//server"), 8)
	assert.Equal(t, getRootLength("//server/share"), 9)
	assert.Equal(t, getRootLength("\\\\server"), 8)
	assert.Equal(t, getRootLength("\\\\server\\share"), 9)
	assert.Equal(t, getRootLength("file:///"), 8)
	assert.Equal(t, getRootLength("file:///path"), 8)
	assert.Equal(t, getRootLength("file:///c:"), 10)
	assert.Equal(t, getRootLength("file:///c:d"), 8)
	assert.Equal(t, getRootLength("file:///c:/path"), 11)
	assert.Equal(t, getRootLength("file:///c%3a"), 12)
	assert.Equal(t, getRootLength("file:///c%3ad"), 8)
	assert.Equal(t, getRootLength("file:///c%3a/path"), 13)
	assert.Equal(t, getRootLength("file:///c%3A"), 12)
	assert.Equal(t, getRootLength("file:///c%3Ad"), 8)
	assert.Equal(t, getRootLength("file:///c%3A/path"), 13)
	assert.Equal(t, getRootLength("file://localhost"), 16)
	assert.Equal(t, getRootLength("file://localhost/"), 17)
	assert.Equal(t, getRootLength("file://localhost/path"), 17)
	assert.Equal(t, getRootLength("file://localhost/c:"), 19)
	assert.Equal(t, getRootLength("file://localhost/c:d"), 17)
	assert.Equal(t, getRootLength("file://localhost/c:/path"), 20)
	assert.Equal(t, getRootLength("file://localhost/c%3a"), 21)
	assert.Equal(t, getRootLength("file://localhost/c%3ad"), 17)
	assert.Equal(t, getRootLength("file://localhost/c%3a/path"), 22)
	assert.Equal(t, getRootLength("file://localhost/c%3A"), 21)
	assert.Equal(t, getRootLength("file://localhost/c%3Ad"), 17)
	assert.Equal(t, getRootLength("file://localhost/c%3A/path"), 22)
	assert.Equal(t, getRootLength("file://server"), 13)
	assert.Equal(t, getRootLength("file://server/"), 14)
	assert.Equal(t, getRootLength("file://server/path"), 14)
	assert.Equal(t, getRootLength("file://server/c:"), 14)
	assert.Equal(t, getRootLength("file://server/c:d"), 14)
	assert.Equal(t, getRootLength("file://server/c:/d"), 14)
	assert.Equal(t, getRootLength("file://server/c%3a"), 14)
	assert.Equal(t, getRootLength("file://server/c%3ad"), 14)
	assert.Equal(t, getRootLength("file://server/c%3a/d"), 14)
	assert.Equal(t, getRootLength("file://server/c%3A"), 14)
	assert.Equal(t, getRootLength("file://server/c%3Ad"), 14)
	assert.Equal(t, getRootLength("file://server/c%3A/d"), 14)
	assert.Equal(t, getRootLength("http://server"), 13)
	assert.Equal(t, getRootLength("http://server/path"), 14)
}

func TestPathIsAbsolute(t *testing.T) {
	// POSIX
	assert.Equal(t, pathIsAbsolute("/path/to/file.ext"), true)
	// DOS
	assert.Equal(t, pathIsAbsolute("c:/path/to/file.ext"), true)
	// URL
	assert.Equal(t, pathIsAbsolute("file:///path/to/file.ext"), true)
	// Non-absolute
	assert.Equal(t, pathIsAbsolute("path/to/file.ext"), false)
	assert.Equal(t, pathIsAbsolute("./path/to/file.ext"), false)
}
