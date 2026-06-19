package luchta

import (
	"os"
	"path/filepath"
	"strings"
)

// RelativizeOutputs makes absolute output paths relative to cwd (forward slashes).
func RelativizeOutputs(cwd string, outputs []string) []string {
	out := make([]string, 0, len(outputs))
	for _, o := range outputs {
		rel, err := filepath.Rel(cwd, o)
		if err != nil {
			rel = o
		}
		out = append(out, filepath.ToSlash(rel))
	}
	return out
}

// CleanOutputs removes *.d.ts and *.d.ts.map files under outDir whose originating
// source file under rootDir no longer exists. No-op when noEmit is true or outDir
// is absent.
func CleanOutputs(cwd, outDir, rootDir string, noEmit bool) error {
	if noEmit {
		return nil
	}
	absOut := filepath.Join(cwd, outDir)
	absRoot := filepath.Join(cwd, rootDir)
	if _, err := os.Stat(absOut); os.IsNotExist(err) {
		return nil
	}
	return filepath.WalkDir(absOut, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		base := d.Name()
		var stem string
		switch {
		case strings.HasSuffix(base, ".d.ts.map"):
			stem = strings.TrimSuffix(base, ".d.ts.map")
		case strings.HasSuffix(base, ".d.ts"):
			stem = strings.TrimSuffix(base, ".d.ts")
		default:
			return nil
		}
		rel, err := filepath.Rel(absOut, filepath.Join(filepath.Dir(path), stem))
		if err != nil {
			return nil
		}
		// source could be .ts or .tsx
		if fileExists(filepath.Join(absRoot, rel+".ts")) || fileExists(filepath.Join(absRoot, rel+".tsx")) {
			return nil
		}
		return os.Remove(path)
	})
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
