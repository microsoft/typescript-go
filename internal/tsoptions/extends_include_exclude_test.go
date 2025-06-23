package tsoptions_test

import (
	"os"
	"path/filepath"
	"testing"
	
	"github.com/microsoft/typescript-go/internal/tsoptions"
	"github.com/microsoft/typescript-go/internal/tsoptions/tsoptionstest"
	"github.com/microsoft/typescript-go/internal/vfs/osvfs"
)

func TestExtendsIncludeExcludeMerging(t *testing.T) {
	// Test that include and exclude arrays are also properly merged
	tempDir, err := os.MkdirTemp("", "tsconfig_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)
	
	// Create base config with include/exclude
	baseDir := filepath.Join(tempDir, "base")
	os.MkdirAll(baseDir, 0755)
	
	baseConfig := `{
  "include": [
    "src/**/*",
    "types/**/*"
  ],
  "exclude": [
    "**/*.test.ts"
  ],
  "compilerOptions": {
    "target": "es2017"
  }
}`
	
	// Create extending config with additional include/exclude
	extendingDir := filepath.Join(tempDir, "extending")
	os.MkdirAll(extendingDir, 0755)
	
	extendingConfig := `{
  "extends": "../base/tsconfig.json",
  "include": [
    "lib/**/*"
  ],
  "exclude": [
    "**/*.spec.ts"
  ],
  "compilerOptions": {
    "outDir": "dist"
  }
}`
	
	// Write configs
	baseConfigPath := filepath.Join(baseDir, "tsconfig.json")
	extendingConfigPath := filepath.Join(extendingDir, "tsconfig.json")
	
	os.WriteFile(baseConfigPath, []byte(baseConfig), 0644)
	os.WriteFile(extendingConfigPath, []byte(extendingConfig), 0644)
	
	// Create some test files
	os.MkdirAll(filepath.Join(extendingDir, "src"), 0755)
	os.MkdirAll(filepath.Join(extendingDir, "types"), 0755)
	os.MkdirAll(filepath.Join(extendingDir, "lib"), 0755)
	
	os.WriteFile(filepath.Join(extendingDir, "src", "main.ts"), []byte("export {}"), 0644)
	os.WriteFile(filepath.Join(extendingDir, "types", "global.d.ts"), []byte("export {}"), 0644)
	os.WriteFile(filepath.Join(extendingDir, "lib", "util.ts"), []byte("export {}"), 0644)
	os.WriteFile(filepath.Join(extendingDir, "src", "test.test.ts"), []byte("export {}"), 0644)
	os.WriteFile(filepath.Join(extendingDir, "src", "spec.spec.ts"), []byte("export {}"), 0644)
	
	// Parse the extending config
	host := &tsoptionstest.VfsParseConfigHost{
		Vfs:              osvfs.FS(),
		CurrentDirectory: extendingDir,
	}
	
	configPath, _ := tsoptions.ParseConfigFileTextToJson(extendingConfigPath, "", extendingConfig)
	result := tsoptions.ParseJsonConfigFileContent(
		configPath,
		host,
		extendingDir,
		nil,
		extendingConfigPath,
		nil,
		nil,
		nil,
	)
	
	t.Logf("Files found: %v", result.ParsedConfig.FileNames)
	
	// Check if files from merged include patterns are found
	hasMain := false
	hasGlobal := false
	hasUtil := false
	hasTestFile := false // Should be excluded by .test.ts pattern
	hasSpecFile := false // Should be excluded by .spec.ts pattern
	
	for _, file := range result.ParsedConfig.FileNames {
		base := filepath.Base(file)
		switch base {
		case "main.ts":
			hasMain = true
		case "global.d.ts":
			hasGlobal = true
		case "util.ts":
			hasUtil = true
		case "test.test.ts":
			hasTestFile = true
		case "spec.spec.ts":
			hasSpecFile = true
		}
	}
	
	if !hasMain {
		t.Error("main.ts should be included (from base include pattern)")
	}
	if !hasGlobal {
		t.Error("global.d.ts should be included (from base include pattern)")
	}
	if !hasUtil {
		t.Error("util.ts should be included (from extending include pattern)")
	}
	if hasTestFile {
		t.Error("test.test.ts should be excluded (from base exclude pattern)")
	}
	if hasSpecFile {
		t.Error("spec.spec.ts should be excluded (from extending exclude pattern)")
	}
}