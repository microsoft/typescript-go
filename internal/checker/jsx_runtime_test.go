// Test to ensure jsx-runtime imports are allowed in declaration files
package checker

import (
	"testing"
)

func TestJSXRuntimeImport(t *testing.T) {
	tests := []struct {
		name     string
		specifier string
		expected bool
	}{
		{"react jsx-runtime", "react/jsx-runtime", true},
		{"react jsx-dev-runtime", "react/jsx-dev-runtime", true},
		{"types jsx-runtime", "@types/react/jsx-runtime", true},
		{"node_modules jsx-runtime", "node_modules/@types/react/jsx-runtime.js", true},
		{"relative jsx-runtime", "../node_modules/@types/react/jsx-runtime.js", true},
		{"nested jsx-runtime", "some/deep/path/jsx-runtime", true},
		{"regular import", "lodash", false},
		{"other node_modules", "node_modules/lodash/index.js", false},
		{"random path", "some/random/path", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isJSXRuntimeImport(tt.specifier)
			if result != tt.expected {
				t.Errorf("isJSXRuntimeImport(%q) = %v, expected %v", tt.specifier, result, tt.expected)
			}
		})
	}
}