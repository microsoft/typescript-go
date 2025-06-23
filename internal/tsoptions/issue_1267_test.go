package tsoptions_test

import (
	"os"
	"path/filepath"
	"testing"
	
	"github.com/microsoft/typescript-go/internal/tsoptions"
	"github.com/microsoft/typescript-go/internal/tsoptions/tsoptionstest"
	"github.com/microsoft/typescript-go/internal/vfs/osvfs"
)

// TestIssue1267ReproducesExactScenario reproduces the exact scenario described in issue #1267
func TestIssue1267ReproducesExactScenario(t *testing.T) {
	// Create a temporary directory structure to test the exact issue scenario
	tempDir, err := os.MkdirTemp("", "tsconfig_issue_1267")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)
	
	// Create base tsconfig that has a "files" block with type declarations
	baseDir := filepath.Join(tempDir, "tsconfig-base")
	os.MkdirAll(baseDir, 0755)
	
	baseConfig := `{
  "$schema": "https://json.schemastore.org/tsconfig",
  "display": "Backend",
  "compilerOptions": {
    "allowJs": true,
    "module": "nodenext",
    "removeComments": true,
    "emitDecoratorMetadata": true,
    "experimentalDecorators": true,
    "allowSyntheticDefaultImports": true,
    "target": "esnext",
    "lib": ["ESNext"],
    "incremental": false,
    "esModuleInterop": true,
    "noImplicitAny": true,
    "moduleResolution": "nodenext",
    "types": ["node", "vitest/globals"],
    "sourceMap": true,
    "strictPropertyInitialization": false
  },
  "files": [
    "types/ical2json.d.ts",
    "types/express.d.ts",
    "types/multer.d.ts",
    "types/reset.d.ts",
    "types/stripe-custom-typings.d.ts",
    "types/nestjs-modules.d.ts",
    "types/luxon.d.ts",
    "types/nestjs-pino.d.ts"
  ],
  "ts-node": {
    "files": true
  }
}`
	
	// Create extending config (monorepo package)
	packageDir := filepath.Join(tempDir, "package")
	os.MkdirAll(packageDir, 0755)
	
	packageConfig := `{
  "extends": "../tsconfig-base/backend.json",
  "compilerOptions": {
    "baseUrl": "./",
    "outDir": "dist",
    "rootDir": "src",
    "resolveJsonModule": true
  },
  "exclude": ["node_modules", "dist"],
  "include": ["src/**/*"]
}`
	
	// Write configs
	baseConfigPath := filepath.Join(baseDir, "backend.json")
	packageConfigPath := filepath.Join(packageDir, "tsconfig.json")
	
	os.WriteFile(baseConfigPath, []byte(baseConfig), 0644)
	os.WriteFile(packageConfigPath, []byte(packageConfig), 0644)
	
	// Create the type declaration files mentioned in the issue
	typesDir := filepath.Join(baseDir, "types")
	os.MkdirAll(typesDir, 0755)
	
	// Create the luxon.d.ts file that was specifically mentioned in the issue
	luxonContent := `declare module 'luxon' {
  interface TSSettings {
    throwOnInvalid: true
  }
}
export {}`
	os.WriteFile(filepath.Join(typesDir, "luxon.d.ts"), []byte(luxonContent), 0644)
	
	// Create other type files
	typeFiles := []string{"ical2json.d.ts", "express.d.ts", "multer.d.ts", "reset.d.ts", 
		"stripe-custom-typings.d.ts", "nestjs-modules.d.ts", "nestjs-pino.d.ts"}
	for _, typeFile := range typeFiles {
		os.WriteFile(filepath.Join(typesDir, typeFile), []byte("export {}"), 0644)
	}
	
	// Create some source files in the package
	srcDir := filepath.Join(packageDir, "src")
	os.MkdirAll(srcDir, 0755)
	os.WriteFile(filepath.Join(srcDir, "main.ts"), []byte("export {}"), 0644)
	os.WriteFile(filepath.Join(srcDir, "utils.ts"), []byte("export {}"), 0644)
	
	// Parse the package config (this is what tsg would do)
	host := &tsoptionstest.VfsParseConfigHost{
		Vfs:              osvfs.FS(),
		CurrentDirectory: packageDir,
	}
	
	configPath, _ := tsoptions.ParseConfigFileTextToJson(packageConfigPath, "", packageConfig)
	result := tsoptions.ParseJsonConfigFileContent(
		configPath,
		host,
		packageDir,
		nil,
		packageConfigPath,
		nil,
		nil,
		nil,
	)
	
	t.Logf("Files found: %v", result.ParsedConfig.FileNames)
	
	// Check that the luxon.d.ts file is included (this was the main issue)
	hasLuxon := false
	hasOtherTypeFiles := 0
	hasSourceFiles := 0
	
	for _, file := range result.ParsedConfig.FileNames {
		fileName := filepath.Base(file)
		if fileName == "luxon.d.ts" {
			hasLuxon = true
		}
		if filepath.Ext(fileName) == ".ts" && filepath.Dir(file) == filepath.Join(packageDir, "types") {
			hasOtherTypeFiles++
		}
		if filepath.Ext(fileName) == ".ts" && filepath.Dir(file) == filepath.Join(packageDir, "src") {
			hasSourceFiles++
		}
	}
	
	if !hasLuxon {
		t.Error("luxon.d.ts file is missing! This was the main issue described in #1267")
	}
	
	// Should have all 8 type files from the base config
	if hasOtherTypeFiles < 7 { // 7 other type files + luxon = 8 total
		t.Errorf("Expected at least 7 other type files, got %d", hasOtherTypeFiles)
	}
	
	// Should also have source files from include pattern
	if hasSourceFiles < 2 {
		t.Errorf("Expected at least 2 source files, got %d", hasSourceFiles)
	}
	
	t.Logf("Successfully found luxon.d.ts and other type declarations from extended config")
	t.Logf("Found %d type files and %d source files", hasOtherTypeFiles+1, hasSourceFiles)
}