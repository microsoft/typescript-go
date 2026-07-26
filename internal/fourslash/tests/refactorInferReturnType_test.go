package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

const inferReturnTypeTitle = "Infer function return type"

// --- Available cases ---

func TestRefactorInferReturnType_available(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `function simple/*marker*/() {
    return 42;
}`
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToMarker(t, "marker")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title:          inferReturnTypeTitle,
		NewFileContent: `function simple(): number {
    return 42;
}`,
	})
}

func TestRefactorInferReturnType_available_arrowFunction(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `const arrow = /*marker*/(x: number) => x * 2;`
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToMarker(t, "marker")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title:          inferReturnTypeTitle,
		NewFileContent: `const arrow = (x: number): number => x * 2;`,
	})
}

func TestRefactorInferReturnType_available_arrowFunction_parenLess(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `const f = /*marker*/x => x * 2;`
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToMarker(t, "marker")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title:          inferReturnTypeTitle,
		NewFileContent: `const f = (x): number => x * 2;`,
	})
}

func TestRefactorInferReturnType_available_methodDeclaration(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `class MyClass {
    myMethod/*marker*/() {
        return 42;
    }
}`
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToMarker(t, "marker")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title:          inferReturnTypeTitle,
		NewFileContent: `class MyClass {
    myMethod(): number {
        return 42;
    }
}`,
	})
}

func TestRefactorInferReturnType_available_overloads(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `function overload/*1*/(x: number): number;
function overload/*2*/(x: string): string;
function overload/*3*/(x: any) {
    return x;
}`
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToMarker(t, "3")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title:          inferReturnTypeTitle,
		NewFileContent: `function overload(x: number): number;
function overload(x: string): string;
function overload(x: any): string | number {
    return x;
}`,
	})
}

func TestRefactorInferReturnType_available_asyncFunction(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `async function asyncFunc/*marker*/() {
    return 42;
}`
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToMarker(t, "marker")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title:          inferReturnTypeTitle,
		NewFileContent: `async function asyncFunc(): Promise<number> {
    return 42;
}`,
	})
}

func TestRefactorInferReturnType_available_generatorFunction(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `function* generator/*marker*/() {
    yield 1;
    yield 2;
}`
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToMarker(t, "marker")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title:          inferReturnTypeTitle,
		NewFileContent: `function* generator(): Generator<1 | 2, void, unknown> {
    yield 1;
    yield 2;
}`,
	})
}

func TestRefactorInferReturnType_available_genericFunction(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `function generic/*marker*/<T>(x: T) {
    return x;
}`
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToMarker(t, "marker")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title:          inferReturnTypeTitle,
		NewFileContent: `function generic<T>(x: T): T {
    return x;
}`,
	})
}

func TestRefactorInferReturnType_available_methodExpression(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `const obj = {
    myMethod/*marker*/() {
        return 42;
    }
};`
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToMarker(t, "marker")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title:          inferReturnTypeTitle,
		NewFileContent: `const obj = {
    myMethod(): number {
        return 42;
    }
};`,
	})
}

func TestRefactorInferReturnType_available_unionReturn(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `function union/*marker*/(flag: boolean) {
    return flag ? 42 : "hello";
}`
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToMarker(t, "marker")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title:          inferReturnTypeTitle,
		NewFileContent: `function union(flag: boolean): "hello" | 42 {
    return flag ? 42 : "hello";
}`,
	})
}

func TestRefactorInferReturnType_available_complexReturn(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `function complex/*marker*/() {
    return { a: 1, b: "hello", c: [1, 2, 3] };
}`
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToMarker(t, "marker")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title:          inferReturnTypeTitle,
		NewFileContent: `function complex(): {
    a: number;
    b: string;
    c: number[];
} {
    return { a: 1, b: "hello", c: [1, 2, 3] };
}`,
	})
}

func TestRefactorInferReturnType_available_voidReturn(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `function voidFunc/*marker*/() {
    console.log("hello");
}`
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToMarker(t, "marker")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title:          inferReturnTypeTitle,
		NewFileContent: `function voidFunc(): void {
    console.log("hello");
}`,
	})
}

func TestRefactorInferReturnType_available_exportedFunction(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `export function/*marker*/() {
    return 42;
}`
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToMarker(t, "marker")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title:          inferReturnTypeTitle,
		NewFileContent: `export function(): number {
    return 42;
}`,
	})
}

func TestRefactorInferReturnType_available_defaultExportedFunction(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `export default function/*marker*/() {
    return 42;
}`
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToMarker(t, "marker")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title:          inferReturnTypeTitle,
		NewFileContent: `export default function(): number {
    return 42;
}`,
	})
}

func TestRefactorInferReturnType_available_computedProperty(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `const key = "foo";
function withComputed/*marker*/() {
    return { [key]: 42 };
}`
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToMarker(t, "marker")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title:          inferReturnTypeTitle,
		NewFileContent: "const key = \"foo\";\nfunction withComputed(): {\n    foo: number;\n} {\n    return { [key]: 42 };\n}",
	})
}

func TestRefactorInferReturnType_available_parameterDeclaration(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `function func(x/*marker*/) {
    return x;
}`
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToMarker(t, "marker")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title:          inferReturnTypeTitle,
		NewFileContent: `function func(x): any {
    return x;
}`,
	})
}

func TestRefactorInferReturnType_available_typePredicate(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `function isString/*marker*/(x: unknown) {
    return typeof x === "string";
}`
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToMarker(t, "marker")
	f.VerifyRefactor(t, fourslash.VerifyRefactorOptions{
		Title:          inferReturnTypeTitle,
		NewFileContent: `function isString(x: unknown): x is string {
    return typeof x === "string";
}`,
	})
}

// --- Not available cases ---

func TestRefactorInferReturnType_notAvailable_withReturnType(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `function simple/*marker*/(): number {
    return 42;
}`
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToMarker(t, "marker")
	f.VerifyRefactorNotAvailable(t, inferReturnTypeTitle)
}

func TestRefactorInferReturnType_notAvailable_constructor(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `class MyClass {
    constructor/*marker*/(x: number) {
    }
}`
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToMarker(t, "marker")
	f.VerifyRefactorNotAvailable(t, inferReturnTypeTitle)
}

func TestRefactorInferReturnType_notAvailable_getAccessor(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `class MyClass {
    get myProp/*marker*/() {
        return 42;
    }
}`
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToMarker(t, "marker")
	f.VerifyRefactorNotAvailable(t, inferReturnTypeTitle)
}

func TestRefactorInferReturnType_notAvailable_setAccessor(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `class MyClass {
    set myProp/*marker*/(value: number) {
    }
}`
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToMarker(t, "marker")
	f.VerifyRefactorNotAvailable(t, inferReturnTypeTitle)
}

func TestRefactorInferReturnType_notAvailable_arrowFunctionExpression(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `const obj = {
    myMethod/*marker*/: (x: number) => x * 2
};`
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToMarker(t, "marker")
	f.VerifyRefactorNotAvailable(t, inferReturnTypeTitle)
}

func TestRefactorInferReturnType_notAvailable_cursorInsideFunctionBody(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `function outer() {
    const inner/*marker*/ = 42;
    return inner;
}`
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToMarker(t, "marker")
	f.VerifyRefactorNotAvailable(t, inferReturnTypeTitle)
}

func TestRefactorInferReturnType_notAvailable_cursorInsideArrowBody(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `const fn = () => {
    const x/*marker*/ = 42;
    return x;
};`
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToMarker(t, "marker")
	f.VerifyRefactorNotAvailable(t, inferReturnTypeTitle)
}

func TestRefactorInferReturnType_notAvailable_constVariable(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `const myConst/*marker*/ = 42;`
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToMarker(t, "marker")
	f.VerifyRefactorNotAvailable(t, inferReturnTypeTitle)
}

func TestRefactorInferReturnType_notAvailable_letVariable(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `let myLet/*marker*/ = 42;`
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToMarker(t, "marker")
	f.VerifyRefactorNotAvailable(t, inferReturnTypeTitle)
}

func TestRefactorInferReturnType_notAvailable_typePredicate(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `function isString/*marker*/(x: any): x is string {
    return typeof x === "string";
}`
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToMarker(t, "marker")
	f.VerifyRefactorNotAvailable(t, inferReturnTypeTitle)
}

func TestRefactorInferReturnType_notAvailable_assertionSignature(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")

	const content = `function assertString/*marker*/(x: any): asserts x is string {
    if (typeof x !== "string") {
        throw new Error("Not a string");
    }
}`
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()
	f.GoToMarker(t, "marker")
	f.VerifyRefactorNotAvailable(t, inferReturnTypeTitle)
}
