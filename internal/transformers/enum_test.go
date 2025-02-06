package transformers

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/binder"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/printer"
	"github.com/microsoft/typescript-go/internal/testutil/emitutil"
	"github.com/microsoft/typescript-go/internal/testutil/parseutil"
)

func TestEnumTransformer(t *testing.T) {
	t.Parallel()
	data := []struct {
		title  string
		input  string
		output string
		jsx    bool
	}{
		{title: "empty enum", input: "enum E {}", output: `var E;
(function (E) {
}(E || (E = {})));`},

		{title: "simple enum", input: "enum E {A}", output: `var E;
(function (E) {
    E[E["A"] = 0] = "A";
}(E || (E = {})));`},

		{title: "autonumber enum #1", input: "enum E {A,B}", output: `var E;
(function (E) {
    E[E["A"] = 0] = "A";
    E[E["B"] = 1] = "B";
}(E || (E = {})));`},

		{title: "autonumber enum #2", input: "enum E {A = 1,B}", output: `var E;
(function (E) {
    E[E["A"] = 1] = "A";
    E[E["B"] = 2] = "B";
}(E || (E = {})));`},

		{title: "autonumber enum #3", input: "enum E {A = 1,B,C}", output: `var E;
(function (E) {
    E[E["A"] = 1] = "A";
    E[E["B"] = 2] = "B";
    E[E["C"] = 3] = "C";
}(E || (E = {})));`},

		{title: "autonumber enum #4", input: "enum E {A = x,B,C}", output: `var E;
(function (E) {
    var auto;
    E[E["A"] = auto = x] = "A";
    E[E["B"] = ++auto] = "B";
    E[E["C"] = ++auto] = "C";
}(E || (E = {})));`},

		{title: "autonumber enum #5", input: "enum E {A = x,B,C = y}", output: `var E;
(function (E) {
    var C, auto;
    E[E["A"] = auto = x] = "A";
    E[E["B"] = ++auto] = "B";
    E["C"] = C = y;
    if (typeof C !== "string") E[C] = "C";
}(E || (E = {})));`},

		{title: "autonumber enum #6", input: "enum E {A = x,B = y,C = z}", output: `var E;
(function (E) {
    var A, B, C;
    E["A"] = A = x;
    if (typeof A !== "string") E[A] = "A";
    E["B"] = B = y;
    if (typeof B !== "string") E[B] = "B";
    E["C"] = C = z;
    if (typeof C !== "string") E[C] = "C";
}(E || (E = {})));`},

		{title: "autonumber enum #7", input: "enum E {A = 1,B,C,D='x'}", output: `var E;
(function (E) {
    E[E["A"] = 1] = "A";
    E[E["B"] = 2] = "B";
    E[E["C"] = 3] = "C";
    E["D"] = 'x';
}(E || (E = {})));`},

		{title: "autonumber enum #8", input: "enum E {A,B=2,C}", output: `var E;
(function (E) {
    E[E["A"] = 0] = "A";
    E[E["B"] = 2] = "B";
    E[E["C"] = 3] = "C";
}(E || (E = {})));`},

		{title: "autonumber enum #9", input: "enum E {A='x',B=2,C}", output: `var E;
(function (E) {
    E["A"] = 'x';
    E[E["B"] = 2] = "B";
    E[E["C"] = 3] = "C";
}(E || (E = {})));`},

		{title: "autonumber enum #10", input: "enum E {A='x',B=y,C}", output: `var E;
(function (E) {
    var auto;
    E["A"] = 'x';
    E[E["B"] = auto = y] = "B";
    E[E["C"] = ++auto] = "C";
}(E || (E = {})));`},

		{title: "autonumber enum #11", input: "enum E {A='x',B=1,C,D=y,E,F=3,G}", output: `var E;
(function (E) {
    var auto;
    E["A"] = 'x';
    E[E["B"] = 1] = "B";
    E[E["C"] = 2] = "C";
    E[E["D"] = auto = y] = "D";
    E[E["E"] = ++auto] = "E";
    E[E["F"] = 3] = "F";
    E[E["G"] = 4] = "G";
}(E || (E = {})));`},

		{title: "autonumber enum #12", input: "enum E {A=-1,B}", output: `var E;
(function (E) {
    E[E["A"] = -1] = "A";
    E[E["B"] = 0] = "B";
}(E || (E = {})));`},

		{title: "autonumber enum #13", input: "enum E {A='x',B}", output: `var E;
(function (E) {
    E["A"] = 'x';
    E["B"] = void 0;
}(E || (E = {})));`},

		{title: "autonumber enum #14", input: "enum E {A,B,C=A|B,D}", output: `var E;
(function (E) {
    var A, B, auto;
    E[E["A"] = A = 0] = "A";
    E[E["B"] = B = 1] = "B";
    E[E["C"] = auto = A | B] = "C";
    E[E["D"] = ++auto] = "D";
}(E || (E = {})));`},

		{title: "string enum", input: "enum E {A = 'x',B = 'y',C = 'z'}", output: `var E;
(function (E) {
    E["A"] = 'x';
    E["B"] = 'y';
    E["C"] = 'z';
}(E || (E = {})));`},

		{title: "number enum", input: "enum E {A = 0,B = 1,C = 2}", output: `var E;
(function (E) {
    E[E["A"] = 0] = "A";
    E[E["B"] = 1] = "B";
    E[E["C"] = 2] = "C";
}(E || (E = {})));`},

		{title: "enum self reference #1", input: "enum E {A,B=A}", output: `var E;
(function (E) {
    var A, B;
    E[E["A"] = A = 0] = "A";
    E["B"] = B = A;
    if (typeof B !== "string") E[B] = "B";
}(E || (E = {})));`},

		{title: "enum self reference #2", input: "enum E {A=x,B=A}", output: `var E;
(function (E) {
    var A, B;
    E["A"] = A = x;
    if (typeof A !== "string") E[A] = "A";
    E["B"] = B = A;
    if (typeof B !== "string") E[B] = "B";
}(E || (E = {})));`},

		{title: "enum self reference #3", input: "enum E {'A'=x,B=A}", output: `var E;
(function (E) {
    var A, B;
    E["A"] = A = x;
    if (typeof A !== "string") E[A] = "A";
    E["B"] = B = A;
    if (typeof B !== "string") E[B] = "B";
}(E || (E = {})));`},

		{title: "enum self reference #4", input: "enum E {'A'=x,'B '=A}", output: `var E;
(function (E) {
    var A;
    E["A"] = A = x;
    if (typeof A !== "string") E[A] = "A";
    E["B "] = A;
    if (typeof E["B "] !== "string") E[E["B "]] = "B ";
}(E || (E = {})));`},

		{title: "export enum", input: "export enum E {A, B}", output: `export var E;
(function (E) {
    E[E["A"] = 0] = "A";
    E[E["B"] = 1] = "B";
}(E || (E = {})));`},

		{title: "const enum", input: "const enum E {A, B}", output: `var E;
(function (E) {
    E[E["A"] = 0] = "A";
    E[E["B"] = 1] = "B";
}(E || (E = {})));`},

		{title: "merged enum", input: "enum E {A}\nenum E {B=A}", output: `var E;
(function (E) {
    E[E["A"] = 0] = "A";
}(E || (E = {})));
(function (E) {
    var B;
    E["B"] = B = A;
    if (typeof B !== "string") E[B] = "B";
}(E || (E = {})));`},
	}

	for _, rec := range data {
		t.Run(rec.title, func(t *testing.T) {
			t.Parallel()
			options := &core.CompilerOptions{}
			file := parseutil.ParseTypeScript(rec.input, rec.jsx)
			parseutil.CheckDiagnostics(t, file)
			binder.BindSourceFile(file, options)
			emitContext := printer.NewEmitContext()
			emitutil.CheckEmit(t, emitContext, NewEnumTransformer(emitContext, options).VisitSourceFile(file), rec.output)
		})
	}
}
