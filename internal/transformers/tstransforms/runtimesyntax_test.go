package tstransforms_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/binder"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/printer"
	"github.com/microsoft/typescript-go/internal/testutil/emittestutil"
	"github.com/microsoft/typescript-go/internal/testutil/parsetestutil"
	"github.com/microsoft/typescript-go/internal/transformers"
	"github.com/microsoft/typescript-go/internal/transformers/tstransforms"
)

func TestNamespaceTransformer(t *testing.T) {
	t.Parallel()
	data := []struct {
		title  string
		input  string
		output string
	}{
		{title: "empty namespace", input: "namespace N {}", output: ``},

		{title: "export var", input: "namespace N { export var x = 1; }", output: `var N;
(function (N) {
    N.x = 1;
})(N || (N = {}));`},

		{title: "export uninitialized var", input: "namespace N { export var x; }", output: `var N;
(function (N) {
})(N || (N = {}));`},

		{title: "exported var reference", input: "namespace N { export var x = 1; x; }", output: `var N;
(function (N) {
    N.x = 1;
    N.x;
})(N || (N = {}));`},

		{title: "exported var reference across namespaces", input: "namespace N { export var x = 1; } namespace N { x; }", output: `var N;
(function (N) {
    N.x = 1;
})(N || (N = {}));
(function (N) {
    N.x;
})(N || (N = {}));`},

		{title: "exported array binding pattern", input: "namespace N { export var [x] = [1]; }", output: `var N;
(function (N) {
    [N.x] = [1];
})(N || (N = {}));`},

		{title: "exported array binding pattern + initializer", input: "namespace N { export var [x = 2] = [1]; }", output: `var N;
(function (N) {
    [N.x = 2] = [1];
})(N || (N = {}));`},

		{title: "exported array binding pattern + elision", input: "namespace N { export var [, x] = [1]; }", output: `var N;
(function (N) {
    [, N.x] = [1];
})(N || (N = {}));`},

		{title: "exported array binding pattern + rest", input: "namespace N { export var [, ...x] = [1]; }", output: `var N;
(function (N) {
    [, ...N.x] = [1];
})(N || (N = {}));`},

		{title: "exported array binding pattern + nested array pattern", input: "namespace N { export var [[x]] = [[1]]; }", output: `var N;
(function (N) {
    [[N.x]] = [[1]];
})(N || (N = {}));`},

		{title: "exported array binding pattern + nested object pattern", input: "namespace N { export var [{x}] = [{x: 1}]; }", output: `var N;
(function (N) {
    [{ x: N.x }] = [{ x: 1 }];
})(N || (N = {}));`},

		{title: "exported object binding pattern", input: "namespace N { export var {x: x} = {x: 1}; }", output: `var N;
(function (N) {
    ({ x: N.x } = { x: 1 });
})(N || (N = {}));`},

		{title: "exported object binding pattern + shorthand assignment", input: "namespace N { export var {x} = {x: 1}; }", output: `var N;
(function (N) {
    ({ x: N.x } = { x: 1 });
})(N || (N = {}));`},

		{title: "exported object binding pattern + initializer", input: "namespace N { export var {x: x = 2} = {x: 1}; }", output: `var N;
(function (N) {
    ({ x: N.x = 2 } = { x: 1 });
})(N || (N = {}));`},

		{title: "exported object binding pattern + shorthand assignment + initializer", input: "namespace N { export var {x = 2} = {x: 1}; }", output: `var N;
(function (N) {
    ({ x: N.x = 2 } = { x: 1 });
})(N || (N = {}));`},

		{title: "exported object binding pattern + rest", input: "namespace N { export var {...x} = {x: 1}; }", output: `var N;
(function (N) {
    ({ ...N.x } = { x: 1 });
})(N || (N = {}));`},

		{title: "exported object binding pattern + nested object pattern", input: "namespace N { export var {y:{x}} = {y: {x: 1}}; }", output: `var N;
(function (N) {
    ({ y: { x: N.x } } = { y: { x: 1 } });
})(N || (N = {}));`},

		{title: "exported object binding pattern + nested array pattern", input: "namespace N { export var {y:[x]} = {y: [1]}; }", output: `var N;
(function (N) {
    ({ y: [N.x] } = { y: [1] });
})(N || (N = {}));`},

		{title: "export function", input: "namespace N { export function f() {} }", output: `var N;
(function (N) {
    function f() { }
    N.f = f;
})(N || (N = {}));`},

		{title: "export class", input: "namespace N { export class C {} }", output: `var N;
(function (N) {
    class C {
    }
    N.C = C;
})(N || (N = {}));`},

		{title: "export namespace", input: "namespace N { export namespace N2 {} }", output: ``},

		{title: "nested namespace", input: "namespace N.N2 { }", output: ``},

		{title: "import=", input: "import X = Y.X;", output: `var X = Y.X;`},

		{title: "export import= at top-level", input: "export import X = Y.X;", output: `export var X = Y.X;`},

		{title: "export import= in namespace", input: "namespace N { export import X = Y.X; }", output: `var N;
(function (N) {
    N.X = Y.X;
})(N || (N = {}));`},

		{title: "shorthand property assignment", input: "namespace N { export var x = 1; var y = { x }; }", output: `var N;
(function (N) {
    N.x = 1;
    var y = { x: N.x };
})(N || (N = {}));`},

		{title: "shorthand property assignment pattern", input: "namespace N { export var x; ({x} = {x: 1}); }", output: `var N;
(function (N) {
    ({ x: N.x } = { x: 1 });
})(N || (N = {}));`},

		{title: "identifier reference in template", input: `namespace N {
    export var x = 1;
    ` + "`" + `${x}` + "`" + `
}`, output: `var N;
(function (N) {
    N.x = 1;
    ` + "`" + `${N.x}` + "`" + `;
})(N || (N = {}));`},
	}

	for _, rec := range data {
		t.Run(rec.title, func(t *testing.T) {
			t.Parallel()
			options := &core.CompilerOptions{}
			file := parsetestutil.ParseTypeScript(rec.input, false /*jsx*/)
			parsetestutil.CheckDiagnostics(t, file)
			binder.BindSourceFile(file)
			emitContext := printer.NewEmitContext()
			resolver := binder.NewReferenceResolver(options, binder.ReferenceResolverHooks{})
			emittestutil.CheckEmit(t, emitContext, tstransforms.NewRuntimeSyntaxTransformer(&transformers.TransformOptions{CompilerOptions: options, Context: emitContext, Resolver: resolver}).TransformSourceFile(file), rec.output)
		})
	}
}

func TestParameterPropertyTransformer(t *testing.T) {
	t.Parallel()
	data := []struct {
		title  string
		input  string
		output string
	}{
		{title: "parameter properties", input: "class C { constructor(public x) { } }", output: `class C {
    x;
    constructor(x) {
        this.x = x;
    }
}`},
		{title: "parameter properties #2", input: "class C extends B { constructor(public x) { super(); } }", output: `class C extends B {
    x;
    constructor(x) {
        super();
        this.x = x;
    }
}`},
	}

	for _, rec := range data {
		t.Run(rec.title, func(t *testing.T) {
			t.Parallel()
			options := &core.CompilerOptions{}
			file := parsetestutil.ParseTypeScript(rec.input, false /*jsx*/)
			parsetestutil.CheckDiagnostics(t, file)
			binder.BindSourceFile(file)
			emitContext := printer.NewEmitContext()
			resolver := binder.NewReferenceResolver(options, binder.ReferenceResolverHooks{})
			opts := &transformers.TransformOptions{Context: emitContext, CompilerOptions: options, Resolver: resolver}
			file = tstransforms.NewTypeEraserTransformer(opts).TransformSourceFile(file)
			file = tstransforms.NewRuntimeSyntaxTransformer(opts).TransformSourceFile(file)
			emittestutil.CheckEmit(t, emitContext, file, rec.output)
		})
	}
}
