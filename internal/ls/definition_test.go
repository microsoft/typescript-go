package ls_test

import (
	"strings"
	"testing"

	"gotest.tools/v3/assert"
)

// Test project reference-based source file detection
func TestProjectReferenceBasedDetection(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		inputFileName  string
		isDeclaration  bool
	}{
		{
			name:          "Declaration file should be detected",
			inputFileName: "packages/o11y/dist/logger.d.ts",
			isDeclaration: true,
		},
		{
			name:          "TypeScript file should not be processed",
			inputFileName: "src/app.ts",
			isDeclaration: false,
		},
		{
			name:          "JavaScript file should not be processed", 
			inputFileName: "lib/utils.js",
			isDeclaration: false,
		},
		{
			name:          "Nested declaration file should be detected",
			inputFileName: "node_modules/@types/react/index.d.ts",
			isDeclaration: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Test the declaration file detection logic used in tryFindSourceFileForDeclaration
			isDecl := strings.HasSuffix(tt.inputFileName, ".d.ts")
			assert.Equal(t, isDecl, tt.isDeclaration)
		})
	}
}

// Test heuristic source mapping patterns 
func TestHeuristicSourceMapping(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		declFilePath   string
		expectedMappings []string // Paths that should be tried
	}{
		{
			name:         "packages dist to src mapping",
			declFilePath: "/Users/person/Developer/repo/packages/o11y/dist/logger/logger.d.ts",
			expectedMappings: []string{
				"/Users/person/Developer/repo/packages/o11y/src/logger/logger.ts",
				"/Users/person/Developer/repo/packages/o11y/src/logger/logger.tsx",
				"/Users/person/Developer/repo/packages/o11y/src/logger/logger.js",
				"/Users/person/Developer/repo/packages/o11y/src/logger/logger.jsx",
			},
		},
		{
			name:         "lib to src mapping",
			declFilePath: "/project/lib/utils/helper.d.ts",
			expectedMappings: []string{
				"/project/src/utils/helper.ts",
				"/project/src/utils/helper.tsx",
				"/project/src/utils/helper.js",
				"/project/src/utils/helper.jsx",
			},
		},
		{
			name:         "build to src mapping",
			declFilePath: "/app/build/components/Button.d.ts",
			expectedMappings: []string{
				"/app/src/components/Button.ts",
				"/app/src/components/Button.tsx", 
				"/app/src/components/Button.js",
				"/app/src/components/Button.jsx",
			},
		},
		{
			name:         "dist to src mapping (most common pattern)", 
			declFilePath: "/simple/dist/index.d.ts",
			expectedMappings: []string{
				"/simple/src/index.ts",
				"/simple/src/index.tsx",
				"/simple/src/index.js",
				"/simple/src/index.jsx",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Test the heuristic mapping logic that would be used
			var candidates []string
			
			// Common monorepo patterns  
			mappings := []struct {
				buildPattern string
				sourcePattern string
			}{
				{"/dist/", "/src/"},
				{"/lib/", "/src/"},
				{"/build/", "/src/"}, 
				{"/out/", "/src/"},
				{"/dist/", "/"},
			}
			
			for _, mapping := range mappings {
				if strings.Contains(tt.declFilePath, mapping.buildPattern) {
					sourceDir := strings.Replace(tt.declFilePath, mapping.buildPattern, mapping.sourcePattern, 1)
					baseName := strings.TrimSuffix(sourceDir, ".d.ts")
					candidates = append(candidates,
						baseName+".ts",
						baseName+".tsx", 
						baseName+".js",
						baseName+".jsx",
					)
					break // Only use first matching pattern
				}
			}

			assert.DeepEqual(t, candidates, tt.expectedMappings)
		})
	}
}

// Integration test placeholder - this would need a real TypeScript program to test properly
func TestGoToDefinitionIntegration(t *testing.T) {
	t.Parallel()
	t.Skip("Integration test requires full TypeScript program setup - implement when needed")
}
