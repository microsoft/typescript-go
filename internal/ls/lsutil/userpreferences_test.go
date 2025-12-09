package lsutil

import (
	"reflect"
	"testing"

	"github.com/go-json-experiment/json"
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
