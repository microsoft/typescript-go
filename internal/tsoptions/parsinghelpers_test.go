package tsoptions

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/microsoft/typescript-go/internal/core"
)

func TestParseCompilerOptionNoMissingFields(t *testing.T) {
	t.Parallel()
	var missingKeys []string
	for _, field := range reflect.VisibleFields(reflect.TypeFor[core.CompilerOptions]()) {
		func() {
			defer func() {
				if r := recover(); r != nil {
					if strings.Contains(fmt.Sprint(r), "interface conversion") {
						// this is a success case, amazingly
						return
					}
					t.Errorf("unexpected panic: %v", r)
				}
			}()
			keyName := field.Name
			// use the JSON key from the tag, if present
			// e.g. `json:"dog[,anythingelse]"` --> dog
			if jsonTag, ok := field.Tag.Lookup("json"); ok {
				keyName = strings.SplitN(jsonTag, ",", 2)[0]
			}
			var something any
			co := core.CompilerOptions{}
			found := parseCompilerOptions(keyName, something, &co)
			if !found {
				missingKeys = append(missingKeys, keyName)
			}
		}()
	}
	if len(missingKeys) > 0 {
		t.Errorf("The following keys are missing entries in the ParseCompilerOptions"+
			" switch statement:\n%v", missingKeys)
	}
}
