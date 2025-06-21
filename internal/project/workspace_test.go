package project_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/project"
	"gotest.tools/v3/assert"
)

func TestWorkspace(t *testing.T) {
	t.Parallel()

	t.Run("GetProjectRootPath", func(t *testing.T) {
		t.Parallel()
		var ws project.Workspace
		ws.SetRoot("/Coding")
		ws.AddFolder("/Coding/One")
		ws.AddFolder("/Coding/Two")
		ws.AddFolder("/Coding/Two/Nested")
		assert.Equal(t, ws.GetProjectRootPath("/Coding/One/file/path.txt"), "/Coding/One")
		assert.Equal(t, ws.GetProjectRootPath("/Coding/Two/file/path.txt"), "/Coding/Two")
		assert.Equal(t, ws.GetProjectRootPath("/Coding/Two/Nested/file"), "/Coding/Two/Nested")
		assert.Equal(t, ws.GetProjectRootPath("/Coding/Two/Nested/f"), "/Coding/Two/Nested")
		assert.Equal(t, ws.GetProjectRootPath("/Coding/One"), "/Coding/One")
		assert.Equal(t, ws.GetProjectRootPath("/Coding/Two"), "/Coding/Two")
		assert.Equal(t, ws.GetProjectRootPath("/Coding/Two/Nested"), "/Coding/Two/Nested")
		assert.Equal(t, ws.GetProjectRootPath("^/untitled/ts-nul-authority/Untitled-1"), "/Coding")
		assert.Equal(t, ws.GetProjectRootPath("/SomethingElse/file.ts"), "")
	})

	t.Run("Removing from folders", func(t *testing.T) {
		t.Parallel()
		var ws project.Workspace
		ws.SetRoot("/Coding")
		ws.AddFolder("/Coding/One")
		ws.AddFolder("/Coding/Two")
		ws.AddFolder("/Coding/Two/Nested")
		ws.RemoveFolder("/Coding/Two")
		assert.Equal(t, ws.GetProjectRootPath("/Coding/One/file/path.txt"), "/Coding/One")
		assert.Equal(t, ws.GetProjectRootPath("/Coding/Two/file/path.txt"), "/Coding")
		assert.Equal(t, ws.GetProjectRootPath("/Coding/Two/Nested/file"), "/Coding/Two/Nested")
		assert.Equal(t, ws.GetProjectRootPath("/Coding/Two/Nested/f"), "/Coding/Two/Nested")
		assert.Equal(t, ws.GetProjectRootPath("/Coding/One"), "/Coding/One")
		assert.Equal(t, ws.GetProjectRootPath("/Coding/Two"), "/Coding")
		assert.Equal(t, ws.GetProjectRootPath("/Coding/Two/Nested"), "/Coding/Two/Nested")
		assert.Equal(t, ws.GetProjectRootPath("^/untitled/ts-nul-authority/Untitled-1"), "/Coding")
		assert.Equal(t, ws.GetProjectRootPath("/SomethingElse/file.ts"), "")
	})
}
