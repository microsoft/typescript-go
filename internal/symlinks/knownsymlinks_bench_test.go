package symlinks

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/tspath"
)

func BenchmarkPopulateSymlinksFromResolutions(b *testing.B) {
	cache := NewKnownSymlink("/project", true)

	deps := make([]struct{ orig, resolved string }, 50)
	for i := range 50 {
		deps[i].orig = "/project/node_modules/pkg" + string(rune('A'+i)) + "/index.js"
		deps[i].resolved = "/real/pkg" + string(rune('A'+i)) + "/index.js"
	}

	b.ResetTimer()
	for range b.N {
		for _, dep := range deps {
			cache.ProcessResolution(dep.orig, dep.resolved)
		}
	}
}

func BenchmarkSetFile(b *testing.B) {
	cache := NewKnownSymlink("/project", true)
	path := tspath.ToPath("/project/file.ts", "/project", true)

	b.ResetTimer()
	for range b.N {
		cache.SetFile(path, "/real/file.ts")
	}
}

func BenchmarkSetDirectory(b *testing.B) {
	cache := NewKnownSymlink("/project", true)
	symlinkPath := tspath.ToPath("/project/symlink", "/project", true).EnsureTrailingDirectorySeparator()
	realDir := &KnownDirectoryLink{
		Real:     "/real/path/",
		RealPath: tspath.ToPath("/real/path", "/project", true).EnsureTrailingDirectorySeparator(),
	}

	b.ResetTimer()
	for range b.N {
		cache.SetDirectory("/project/symlink", symlinkPath, realDir)
	}
}

func BenchmarkGuessDirectorySymlink(b *testing.B) {
	cache := NewKnownSymlink("/project", true)

	b.ResetTimer()
	for range b.N {
		cache.guessDirectorySymlink(
			"/real/node_modules/package/dist/index.js",
			"/project/symlink/package/dist/index.js",
			"/project",
		)
	}
}

func BenchmarkConcurrentAccess(b *testing.B) {
	cache := NewKnownSymlink("/project", true)

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			path := tspath.ToPath("/project/file"+string(rune('A'+(i%26)))+".ts", "/project", true)
			cache.SetFile(path, "/real/file.ts")
			cache.Files().Load(path)
			i++
		}
	})
}
