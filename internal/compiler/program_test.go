package compiler

import (
	"testing"
	"testing/fstest"

	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/vfs/vfstest"
	"gotest.tools/v3/assert"
)

func TestProgramFileOrdering(t *testing.T) {
	t.Parallel()
	fs := fstest.MapFS{}

	files := map[string]string{
		"c:/dev/src/index.ts":          "/// <reference path='c:/dev/src2/a/5.ts' />\n/// <reference path='c:/dev/src2/a/10.ts' />",
		"c:/dev/src2/a/5.ts":           `/// <reference path="4.ts" />`,
		"c:/dev/src2/a/4.ts":           `/// <reference path="b/3.ts" />`,
		"c:/dev/src2/a/b/3.ts":         `/// <reference path="2.ts" />`,
		"c:/dev/src2/a/b/2.ts":         `/// <reference path="c/1.ts" />`,
		"c:/dev/src2/a/b/c/1.ts":       `console.log("hello");`,
		"c:/dev/src2/a/10.ts":          `/// <reference path="b/c/d/9.ts" />`,
		"c:/dev/src2/a/b/c/d/9.ts":     `/// <reference path="e/8.ts" />`,
		"c:/dev/src2/a/b/c/d/e/8.ts":   `/// <reference path="7.ts" />`,
		"c:/dev/src2/a/b/c/d/e/7.ts":   `/// <reference path="f/6.ts" />`,
		"c:/dev/src2/a/b/c/d/e/f/6.ts": `console.log("world!");`,
	}

	for fileName, contents := range files {
		fs[fileName] = &fstest.MapFile{
			Data: []byte(contents),
		}
	}

	opts := core.CompilerOptions{}

	program := NewProgram(ProgramOptions{
		RootPath:       "c:/dev/src",
		Host:           NewCompilerHost(&opts, "c:/dev/src", vfstest.FromMapFS(fs, true)),
		Options:        &opts,
		SingleThreaded: false,
	})

	actualOrder := []string{}
	for _, file := range program.files {
		actualOrder = append(actualOrder, file.FileName())
	}

	expectedOrder := []string{
		"c:/dev/src2/a/b/c/1.ts",
		"c:/dev/src2/a/b/2.ts",
		"c:/dev/src2/a/b/3.ts",
		"c:/dev/src2/a/4.ts",
		"c:/dev/src2/a/5.ts",
		"c:/dev/src2/a/b/c/d/e/f/6.ts",
		"c:/dev/src2/a/b/c/d/e/7.ts",
		"c:/dev/src2/a/b/c/d/e/8.ts",
		"c:/dev/src2/a/b/c/d/9.ts",
		"c:/dev/src2/a/10.ts",
		"c:/dev/src/index.ts",
	}

	assert.DeepEqual(t, expectedOrder, actualOrder)
}
