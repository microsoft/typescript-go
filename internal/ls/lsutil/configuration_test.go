package lsutil

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/core"
)

func TestUserConfig_GetPreferences(t *testing.T) {
	t.Parallel()

	type expectedPreference int

	const (
		expectedPreferenceTS expectedPreference = iota
		expectedPreferenceJS
	)

	doubleQuotePrefs := UserPreferences{
		QuotePreference: "double",
	}
	singleQuotePrefs := UserPreferences{
		QuotePreference: "single",
	}

	tsDoubleQuoteJsSingleQuoteConfig := &UserConfig{
		ts: doubleQuotePrefs,
		js: singleQuotePrefs,
	}
	tests := []struct {
		name         string
		config       *UserConfig
		activeFile   string
		expectedPref expectedPreference
	}{
		{
			name:         ".ts file returns TS preferences",
			config:       tsDoubleQuoteJsSingleQuoteConfig,
			activeFile:   "file.ts",
			expectedPref: expectedPreferenceTS,
		},
		{
			name:         ".tsx file returns TS preferences",
			config:       tsDoubleQuoteJsSingleQuoteConfig,
			activeFile:   "file.tsx",
			expectedPref: expectedPreferenceTS,
		},
		{
			name:         ".d.ts file returns TS preferences",
			config:       tsDoubleQuoteJsSingleQuoteConfig,
			activeFile:   "file.d.ts",
			expectedPref: expectedPreferenceTS,
		},
		{
			name:         ".mts file returns TS preferences",
			config:       tsDoubleQuoteJsSingleQuoteConfig,
			activeFile:   "file.mts",
			expectedPref: expectedPreferenceTS,
		},
		{
			name:         ".cts file returns TS preferences",
			config:       tsDoubleQuoteJsSingleQuoteConfig,
			activeFile:   "file.cts",
			expectedPref: expectedPreferenceTS,
		},
		{
			name:         ".js file returns JS preferences",
			config:       tsDoubleQuoteJsSingleQuoteConfig,
			activeFile:   "file.js",
			expectedPref: expectedPreferenceJS,
		},
		{
			name:         ".jsx file returns JS preferences",
			config:       tsDoubleQuoteJsSingleQuoteConfig,
			activeFile:   "file.jsx",
			expectedPref: expectedPreferenceJS,
		},
		{
			name:         ".mjs file returns JS preferences",
			config:       tsDoubleQuoteJsSingleQuoteConfig,
			activeFile:   "file.mjs",
			expectedPref: expectedPreferenceJS,
		},
		{
			name:         ".cjs file returns JS preferences",
			config:       tsDoubleQuoteJsSingleQuoteConfig,
			activeFile:   "file.cjs",
			expectedPref: expectedPreferenceJS,
		},
		{
			name:         "Empty file returns TS preferences",
			config:       tsDoubleQuoteJsSingleQuoteConfig,
			activeFile:   "",
			expectedPref: expectedPreferenceTS,
		},
		{
			name:         "Unknown file extension returns JS preferences",
			config:       tsDoubleQuoteJsSingleQuoteConfig,
			activeFile:   "file.py",
			expectedPref: expectedPreferenceJS,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.config.GetPreferences(tt.activeFile)

			switch tt.expectedPref {
			case expectedPreferenceTS:
				if got.QuotePreference != doubleQuotePrefs.QuotePreference {
					t.Errorf("GetPreferences(%q) = %v, expected TS preferences (%v)", tt.activeFile, got.QuotePreference, doubleQuotePrefs.QuotePreference)
				}
			case expectedPreferenceJS:
				if got.QuotePreference != singleQuotePrefs.QuotePreference {
					t.Errorf("GetPreferences(%q) = %v, expected JS preferences (%v)", tt.activeFile, got.QuotePreference, singleQuotePrefs.QuotePreference)
				}
			}
		})
	}
}

func TestUserConfig_GetPreferences_CodeLensAndInlayHints(t *testing.T) {
	t.Parallel()

	codeLensAndInlayHintsOn := UserPreferences{
		CodeLens: CodeLensUserPreferences{
			ReferencesCodeLensEnabled: core.TSTrue,
		},
		InlayHints: InlayHintsPreferences{
			IncludeInlayVariableTypeHints: core.TSTrue,
		},
	}
	codeLensAndInlayHintsOff := UserPreferences{
		CodeLens: CodeLensUserPreferences{
			ReferencesCodeLensEnabled: core.TSFalse,
		},
		InlayHints: InlayHintsPreferences{
			IncludeInlayVariableTypeHints: core.TSFalse,
		},
	}

	tests := []struct {
		name                   string
		config                 *UserConfig
		activeFile             string
		expectedLensesAndHints core.Tristate
	}{
		{
			name: ".ts file with CodeLens and InlayHints enabled",
			config: &UserConfig{
				ts: codeLensAndInlayHintsOn,
				js: codeLensAndInlayHintsOff,
			},
			activeFile:             "file.ts",
			expectedLensesAndHints: core.TSTrue,
		},
		{
			name: ".js file with CodeLens and InlayHints disabled",
			config: &UserConfig{
				ts: codeLensAndInlayHintsOn,
				js: codeLensAndInlayHintsOff,
			},
			activeFile:             "file.js",
			expectedLensesAndHints: core.TSFalse,
		},
		{
			name: ".ts file with CodeLens and InlayHints disabled",
			config: &UserConfig{
				ts: codeLensAndInlayHintsOff,
				js: codeLensAndInlayHintsOn,
			},
			activeFile:             "file.ts",
			expectedLensesAndHints: core.TSFalse,
		},
		{
			name: ".js file with CodeLens and InlayHints enabled",
			config: &UserConfig{
				ts: codeLensAndInlayHintsOff,
				js: codeLensAndInlayHintsOn,
			},
			activeFile:             "file.js",
			expectedLensesAndHints: core.TSTrue,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.config.GetPreferences(tt.activeFile)

			if got.CodeLens.ReferencesCodeLensEnabled != tt.expectedLensesAndHints {
				t.Errorf("GetPreferences().CodeLens.ReferencesCodeLensEnabled was '%v', expected '%v'", got.CodeLens.ReferencesCodeLensEnabled, tt.expectedLensesAndHints)
			}
			if got.InlayHints.IncludeInlayVariableTypeHints != tt.expectedLensesAndHints {
				t.Errorf("GetPreferences().InlayHints.IncludeInlayVariableTypeHints was '%v', expected '%v'", got.InlayHints.IncludeInlayVariableTypeHints, tt.expectedLensesAndHints)
			}
		})
	}
}
