package lsutil

import (
	"reflect"
	"testing"

	"github.com/go-json-experiment/json"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/modulespecifiers"
	"gotest.tools/v3/assert"
)

func fillNonZeroValues(v reflect.Value) {
	t := v.Type()
	for i := range t.NumField() {
		field := v.Field(i)
		if !field.CanSet() {
			continue
		}
		switch field.Kind() {
		case reflect.Bool:
			field.SetBool(true)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			field.SetInt(1)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			field.SetUint(1)
		case reflect.String:
			val := getValidStringValue(field.Type())
			field.SetString(val)
		case reflect.Slice:
			if field.Type().Elem().Kind() == reflect.String {
				field.Set(reflect.ValueOf([]string{"test"}))
			}
		case reflect.Struct:
			fillNonZeroValues(field)
		}
	}
}

func getValidStringValue(t reflect.Type) string {
	typeName := t.String()
	switch typeName {
	case "lsutil.QuotePreference":
		return string(QuotePreferenceSingle)
	case "lsutil.JsxAttributeCompletionStyle":
		return string(JsxAttributeCompletionStyleBraces)
	case "lsutil.IncludePackageJsonAutoImports":
		return string(IncludePackageJsonAutoImportsOn)
	case "lsutil.IncludeInlayParameterNameHints":
		return string(IncludeInlayParameterNameHintsAll)
	case "modulespecifiers.ImportModuleSpecifierPreference":
		return string(modulespecifiers.ImportModuleSpecifierPreferenceRelative)
	case "modulespecifiers.ImportModuleSpecifierEndingPreference":
		return string(modulespecifiers.ImportModuleSpecifierEndingPreferenceJs)
	default:
		return "test"
	}
}

func TestUserPreferencesRoundtrip(t *testing.T) {
	t.Parallel()

	original := &UserPreferences{}
	fillNonZeroValues(reflect.ValueOf(original).Elem())

	jsonBytes, err := json.Marshal(original)
	assert.NilError(t, err)

	t.Run("UnmarshalJSONFrom", func(t *testing.T) {
		t.Parallel()
		parsed := &UserPreferences{}
		err2 := json.Unmarshal(jsonBytes, parsed)
		assert.NilError(t, err2)
		assert.DeepEqual(t, original, parsed)
	})

	t.Run("parseWorker", func(t *testing.T) {
		t.Parallel()
		var config map[string]any
		err2 := json.Unmarshal(jsonBytes, &config)
		assert.NilError(t, err2)
		parsed := &UserPreferences{}
		parsed.parseWorker(config)
		assert.DeepEqual(t, original, parsed)
	})
}

func TestUserPreferencesSerialize(t *testing.T) {
	t.Parallel()

	t.Run("config path field serializes to nested path", func(t *testing.T) {
		t.Parallel()
		prefs := &UserPreferences{
			QuotePreference: QuotePreferenceSingle,
		}
		jsonBytes, err := json.Marshal(prefs)
		assert.NilError(t, err)

		var actual map[string]any
		err = json.Unmarshal(jsonBytes, &actual)
		assert.NilError(t, err)

		preferences := actual["preferences"].(map[string]any)
		assert.Equal(t, "single", preferences["quoteStyle"])
	})

	t.Run("raw-only field serializes to unstable section", func(t *testing.T) {
		t.Parallel()
		prefs := &UserPreferences{
			DisableSuggestions: true,
		}
		jsonBytes, err := json.Marshal(prefs)
		assert.NilError(t, err)

		var actual map[string]any
		err = json.Unmarshal(jsonBytes, &actual)
		assert.NilError(t, err)

		unstable := actual["unstable"].(map[string]any)
		assert.Equal(t, true, unstable["disableSuggestions"])
	})

	t.Run("inlay hint inversion on serialize", func(t *testing.T) {
		t.Parallel()
		prefs := &UserPreferences{
			InlayHints: InlayHintsPreferences{
				IncludeInlayParameterNameHints:                        IncludeInlayParameterNameHintsAll,
				IncludeInlayParameterNameHintsWhenArgumentMatchesName: true,
			},
		}
		jsonBytes, err := json.Marshal(prefs)
		assert.NilError(t, err)

		var actual map[string]any
		err = json.Unmarshal(jsonBytes, &actual)
		assert.NilError(t, err)

		inlayHints := actual["inlayHints"].(map[string]any)
		parameterNames := inlayHints["parameterNames"].(map[string]any)
		assert.Equal(t, "all", parameterNames["enabled"])
		assert.Equal(t, false, parameterNames["suppressWhenArgumentMatchesName"]) // inverted
	})

	t.Run("mixed config and unstable fields", func(t *testing.T) {
		t.Parallel()
		prefs := &UserPreferences{
			QuotePreference:      QuotePreferenceSingle,
			DisableSuggestions:   true,
			DisplayPartsForJSDoc: true,
		}
		jsonBytes, err := json.Marshal(prefs)
		assert.NilError(t, err)

		var actual map[string]any
		err = json.Unmarshal(jsonBytes, &actual)
		assert.NilError(t, err)

		preferences := actual["preferences"].(map[string]any)
		assert.Equal(t, "single", preferences["quoteStyle"])

		unstable := actual["unstable"].(map[string]any)
		assert.Equal(t, true, unstable["disableSuggestions"])
		assert.Equal(t, true, unstable["displayPartsForJSDoc"])
	})
}

func TestUserPreferencesParseUnstable(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		json     string
		expected *UserPreferences
	}{
		{
			name: "unstable fields with correct casing",
			json: `{
				"unstable": {
					"disableSuggestions": true,
					"maximumHoverLength": 100,
					"allowRenameOfImportPath": true
				}
			}`,
			expected: &UserPreferences{
				DisableSuggestions:      true,
				MaximumHoverLength:      100,
				AllowRenameOfImportPath: true,
			},
		},
		{
			name: "nested preferences path",
			json: `{
				"preferences": {
					"quoteStyle": "single",
					"useAliasesForRenames": true
				}
			}`,
			expected: &UserPreferences{
				QuotePreference:     QuotePreferenceSingle,
				UseAliasesForRename: core.TSTrue,
			},
		},
		{
			name: "suggest section",
			json: `{
				"suggest": {
					"autoImports": false,
					"includeCompletionsForImportStatements": true
				}
			}`,
			expected: &UserPreferences{
				IncludeCompletionsForModuleExports:    core.TSFalse,
				IncludeCompletionsForImportStatements: core.TSTrue,
			},
		},
		{
			name: "inlayHints with invert",
			json: `{
				"inlayHints": {
					"parameterNames": {
						"enabled": "all",
						"suppressWhenArgumentMatchesName": true
					}
				}
			}`,
			expected: &UserPreferences{
				InlayHints: InlayHintsPreferences{
					IncludeInlayParameterNameHints:                        IncludeInlayParameterNameHintsAll,
					IncludeInlayParameterNameHintsWhenArgumentMatchesName: false, // inverted
				},
			},
		},
		{
			name: "mixed config",
			json: `{
				"unstable": {
					"displayPartsForJSDoc": true
				},
				"preferences": {
					"importModuleSpecifier": "relative"
				},
				"workspaceSymbols": {
					"excludeLibrarySymbols": true
				}
			}`,
			expected: &UserPreferences{
				DisplayPartsForJSDoc: true,
				ModuleSpecifier: ModuleSpecifierUserPreferences{
					ImportModuleSpecifierPreference: modulespecifiers.ImportModuleSpecifierPreferenceRelative,
				},
				ExcludeLibrarySymbolsInNavTo: true,
			},
		},
		{
			name: "stable config overrides unstable",
			json: `{
				"unstable": {
					"quotePreference": "double"
				},
				"preferences": {
					"quoteStyle": "single"
				}
			}`,
			expected: &UserPreferences{
				QuotePreference: QuotePreferenceSingle, // stable wins
			},
		},
		{
			name: "unstable sets value when no stable config",
			json: `{
				"unstable": {
					"includeCompletionsWithSnippetText": false
				}
			}`,
			expected: &UserPreferences{
				IncludeCompletionsWithSnippetText: core.TSFalse,
			},
		},
		{
			name: "any field can be passed via unstable by its raw name",
			json: `{
				"unstable": {
					"quotePreference": "double",
					"includeCompletionsForModuleExports": true,
					"excludeLibrarySymbolsInNavTo": true
				}
			}`,
			expected: &UserPreferences{
				QuotePreference:                    QuotePreferenceDouble,
				IncludeCompletionsForModuleExports: core.TSTrue,
				ExcludeLibrarySymbolsInNavTo:       true,
			},
		},
		{
			name: "TypeScript raw names work in unstable section",
			json: `{
				"unstable": {
					"includeCompletionsForModuleExports": true,
					"quotePreference": "single",
					"providePrefixAndSuffixTextForRename": true,
					"includeInlayParameterNameHints": "all",
					"organizeImportsLocale": "en"
				}
			}`,
			expected: &UserPreferences{
				IncludeCompletionsForModuleExports: core.TSTrue,
				QuotePreference:                    QuotePreferenceSingle,
				UseAliasesForRename:                core.TSTrue,
				OrganizeImportsLocale:              "en",
				InlayHints: InlayHintsPreferences{
					IncludeInlayParameterNameHints: IncludeInlayParameterNameHintsAll,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var config map[string]any
			err := json.Unmarshal([]byte(tt.json), &config)
			assert.NilError(t, err)

			parsed := &UserPreferences{}
			parsed.parseWorker(config)

			assert.DeepEqual(t, tt.expected, parsed)
		})
	}
}
