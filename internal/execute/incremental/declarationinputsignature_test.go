package incremental

import "testing"

func TestDeclarationInputSignature(t *testing.T) {
	t.Parallel()

	base := "export const value = 1;\n"
	baseSignature := computeDeclarationInputSignature(base)
	for _, text := range []string{
		"export  const value = 1;\n",
		"export const value = 1; // comment\n",
		"export const /* comment */ value = 1;\n",
		"/* comment */\nexport const value = 1;\n",
	} {
		if signature := computeDeclarationInputSignature(text); signature != baseSignature {
			t.Errorf("ordinary trivia changed declaration input signature for %q", text)
		}
	}

	for _, text := range []string{
		"export const value = 2;\n",
		"/** @internal */\nexport const value = 1;\n",
		"// @jsxImportSource ./jsx-runtime\nexport const value = 1;\n",
		"//@jsxImportSource ./jsx-runtime\nexport const value = 1;\n",
		"/// <reference path=\"./types.d.ts\" />\nexport const value = 1;\n",
		"export const value = 1;\nexport const other = 2;\n",
	} {
		if signature := computeDeclarationInputSignature(text); signature == baseSignature {
			t.Errorf("declaration input change was ignored for %q", text)
		}
	}

	if computeDeclarationInputSignature("function f() { return\nvalue }\n") ==
		computeDeclarationInputSignature("function f() { return value }\n") {
		t.Error("line break affecting automatic semicolon insertion was ignored")
	}
}
