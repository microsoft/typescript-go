package main

import (
	"context"
	"fmt"
	"os"

	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/ls"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
)

func main() {
	// Read the test file
	content, err := os.ReadFile("test_chinese.ts")
	if err != nil {
		panic(err)
	}

	// Create a simple compilation host 
	host := &compiler.CompilerHost{}
	host.ReadFile = func(fileName string) (string, error) {
		if fileName == "test_chinese.ts" {
			return string(content), nil
		}
		return "", fmt.Errorf("file not found: %s", fileName)
	}

	host.FileExists = func(fileName string) bool {
		return fileName == "test_chinese.ts"
	}

	// Create program
	options := &core.CompilerOptions{
		Target:           core.ScriptTargetES5,
		Module:           core.ModuleKindCommonJS,
		Lib:              []string{"lib.es5.d.ts"},
		AllowJs:          core.TSFalse,
		CheckJs:          core.TSFalse,
		Strict:           core.TSTrue,
		ModuleResolution: core.ModuleResolutionKindNodeJs,
	}

	program := compiler.CreateProgram([]string{"test_chinese.ts"}, options, host, nil, nil)
	
	// Create language service
	lsHost := &ls.LanguageServiceHost{}
	lsHost.FileExists = host.FileExists
	lsHost.ReadFile = func(fileName string) (string, error) {
		if fileName == "test_chinese.ts" {
			return string(content), nil
		}
		return "", fmt.Errorf("file not found: %s", fileName)
	}

	service := ls.CreateLanguageService(lsHost, nil, nil)

	// Try to get completions at the end of the file
	// Position should be after the opening parenthesis in setLengthTextPositionPreset(
	position := lsproto.Position{Line: 12, Character: 47} // After the opening parenthesis
	
	completions, err := service.ProvideCompletion(
		context.Background(),
		lsproto.DocumentUri("file:///test_chinese.ts"),
		position,
		nil,
		nil,
		&ls.UserPreferences{},
	)

	if err != nil {
		panic(err)
	}

	fmt.Printf("Completions found: %d\n", len(completions.List.Items))
	for i, item := range completions.List.Items {
		fmt.Printf("  [%d] %s\n", i, item.Label)
	}
}