package fourslash_test

import (
	"strings"
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestCompletionListInUnclosedTypeArguments(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `let x = 10;
type Type = void;
declare function f<T>(): void;
declare function f2<T, U>(): void;
f<Ty
f<Ty;
f<Ty>
f<Ty>
f<Ty>();

f2<Ty,
f2<Ty,Ty
f2<Ty,Ty;
f2<Ty,Ty>
f2<Ty,Ty>
f2<Ty,{| "newId": true, "typeOnly": true |}Ty>();

f2<typeof x, Ty

f2<Ty, () =>Ty
f2<() =>Ty, () =>Ty
f2<any, () =>Ty`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.GoToEachMarker(t, nil, func(t *testing.T, marker *fourslash.Marker, index int) {
		markerName := marker.Name
		valueOnly := markerName != nil && strings.HasSuffix(*markerName, "ValueOnly")
		commitCharacters := &defaultCommitCharacters
		if marker.Data != nil {
			if newId, ok := marker.Data["newId"]; ok && newId.(bool) {
				commitCharacters = &[]string{".", ";"}
			}
		}
		var includes []fourslash.CompletionsExpectedItem
		var excludes []string
		if valueOnly {
			includes = []fourslash.CompletionsExpectedItem{
				"x",
			}
			excludes = []string{
				"Type",
			}
		} else {
			includes = []fourslash.CompletionsExpectedItem{
				"Type",
			}
			excludes = []string{
				"x",
			}
		}
		f.VerifyCompletions(t, marker, &fourslash.CompletionsExpectedList{
			IsIncomplete: false,
			ItemDefaults: &fourslash.CompletionsExpectedItemDefaults{
				CommitCharacters: commitCharacters,
				EditRange:        ignored,
			},
			Items: &fourslash.CompletionsExpectedItems{
				Includes: includes,
				Excludes: excludes,
			},
		})
	})
}
