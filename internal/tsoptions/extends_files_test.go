package tsoptions_test

import (
	"os"
	"path/filepath"
	"testing"
	
	"github.com/microsoft/typescript-go/internal/tsoptions"
	"github.com/microsoft/typescript-go/internal/tsoptions/tsoptionstest"
	"github.com/microsoft/typescript-go/internal/vfs/osvfs"
)

func TestExtendsFilesMerging(t *testing.T) {
	// Create a temporary directory structure to test
	tempDir, err := os.MkdirTemp("", "tsconfig_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)
	
	// Create base config with files
	baseDir := filepath.Join(tempDir, "base")
	os.MkdirAll(baseDir, 0755)
	
	baseConfig := `{
  "files": [
    "types/luxon.d.ts",
    "types/express.d.ts"
  ],
  "compilerOptions": {
    "target": "es2017"
  }
}`
	
	// Create extending config with its own files
	extendingDir := filepath.Join(tempDir, "extending")
	os.MkdirAll(extendingDir, 0755)
	
	extendingConfig := `{
  "extends": "../base/tsconfig.json",
  "files": [
    "src/main.ts"
  ],
  "compilerOptions": {
    "outDir": "dist"
  }
}`
	
	// Write configs
	baseConfigPath := filepath.Join(baseDir, "tsconfig.json")
	extendingConfigPath := filepath.Join(extendingDir, "tsconfig.json")
	
	err = os.WriteFile(baseConfigPath, []byte(baseConfig), 0644)
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(extendingConfigPath, []byte(extendingConfig), 0644)
	if err != nil {
		t.Fatal(err)
	}
	
	// Create dummy files 
	os.MkdirAll(filepath.Join(baseDir, "types"), 0755)
	os.WriteFile(filepath.Join(baseDir, "types", "luxon.d.ts"), []byte("export {}"), 0644)
	os.WriteFile(filepath.Join(baseDir, "types", "express.d.ts"), []byte("export {}"), 0644)
	
	os.MkdirAll(filepath.Join(extendingDir, "src"), 0755)
	os.WriteFile(filepath.Join(extendingDir, "src", "main.ts"), []byte("export {}"), 0644)
	
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
		nil, // existingOptions
		extendingConfigPath,
		nil, // resolutionStack
		nil, // extraFileExtensions
		nil, // extendedConfigCache
	)
	
	t.Logf("Files found: %v", result.ParsedConfig.FileNames)
	
	// Check if files from base config are included
	hasLuxon := false
	hasExpress := false
	hasMain := false
	
	for _, file := range result.ParsedConfig.FileNames {
		if filepath.Base(file) == "luxon.d.ts" {
			hasLuxon = true
		}
		if filepath.Base(file) == "express.d.ts" {
			hasExpress = true
		}
		if filepath.Base(file) == "main.ts" {
			hasMain = true
		}
	}
	
	t.Logf("Has luxon.d.ts: %v", hasLuxon)
	t.Logf("Has express.d.ts: %v", hasExpress)
	t.Logf("Has main.ts: %v", hasMain)
	
	if !hasLuxon || !hasExpress {
		t.Error("Files from extended config are missing! This indicates the bug exists.")
	}
	if !hasMain {
		t.Error("Files from extending config are missing!")  
	}
}