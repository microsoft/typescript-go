package ls

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/ls"
	"github.com/microsoft/typescript-go/internal/parser"
	"github.com/microsoft/typescript-go/internal/testutil/projecttestutil"
	"github.com/microsoft/typescript-go/internal/tspath"
	"gotest.tools/v3/assert"

	"github.com/microsoft/typescript-go/internal/scanner"
)

type verifySignatureHelpOptions struct {
	//overloadsCount            int
	docComment                string
	text                      string
	parameterName             string
	parameterSpan             string
	parameterDocComment       string
	parameterCount            int
	isVariadic                bool
	triggerReason             ls.SignatureHelpTriggerReason
	overrideSelectedItemIndex int
	//tags?: ReadonlyArray<JSDocTagInfo>;
}

func TestSignature(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		title  string
		input  string
		output map[string]verifySignatureHelpOptions
	}{
		{
			title: "SignatureHelpCallExpressions",
			input: `function fnTest(str: string, num: number) { }
fnTest(/*1*/'', /*2*/5);`,
			output: map[string]verifySignatureHelpOptions{
				"1": {
					text:           `fnTest(str: string, num: number): void`,
					parameterCount: 2,
					parameterSpan:  "str: string",
				},
				"2": {
					text:           `fnTest(str: string, num: number): void`,
					parameterCount: 2,
					parameterSpan:  "num: number",
				},
			},
		},
		{
			title: "SignatureHelp_contextual",
			input: `interface I {
	m(n: number, s: string): void;
	m2: () => void;
}
declare function takesObj(i: I): void;
takesObj({ m: (/*takesObj0*/) });
takesObj({ m(/*takesObj1*/) });
takesObj({ m: function(/*takesObj2*/) });
takesObj({ m2: (/*takesObj3*/) })
declare function takesCb(cb: (n: number, s: string, b: boolean) => void): void;
takesCb((/*contextualParameter1*/));
takesCb((/*contextualParameter1b*/) => {});
takesCb((n, /*contextualParameter2*/));
takesCb((n, s, /*contextualParameter3*/));
takesCb((n,/*contextualParameter3_2*/ s, b));
takesCb((n, s, b, /*contextualParameter4*/))
type Cb = () => void;
const cb: Cb = (/*contextualTypeAlias*/
const cb2: () => void = (/*contextualFunctionType*/)`,
			output: map[string]verifySignatureHelpOptions{
				"takesObj0": {
					text:           "m(n: number, s: string): void",
					parameterCount: 2,
					parameterSpan:  "n: number",
				},
				"takesObj1": {
					text:           "m(n: number, s: string): void",
					parameterCount: 2,
					parameterSpan:  "n: number",
				},
				"takesObj2": {
					text:           "m(n: number, s: string): void",
					parameterCount: 2,
					parameterSpan:  "n: number",
				},
				"takesObj3": {
					text:           "m2(): void",
					parameterCount: 0,
					parameterSpan:  "",
				},
				"contextualParameter1": {
					text:           "cb(n: number, s: string, b: boolean): void",
					parameterCount: 3,
					parameterSpan:  "n: number",
				},
				"contextualParameter1b": {
					text:           "cb(n: number, s: string, b: boolean): void",
					parameterCount: 3,
					parameterSpan:  "n: number",
				},
				"contextualParameter2": {
					text:           "cb(n: number, s: string, b: boolean): void",
					parameterCount: 3,
					parameterSpan:  "s: string",
				},
				"contextualParameter3": {
					text:           "cb(n: number, s: string, b: boolean): void",
					parameterCount: 3,
					parameterSpan:  "b: boolean",
				},
				"contextualParameter3_2": {
					text:           "cb(n: number, s: string, b: boolean): void",
					parameterCount: 3,
					parameterSpan:  "s: string",
				},
				"contextualParameter4": {
					text:           "cb(n: number, s: string, b: boolean): void",
					parameterCount: 3,
					parameterSpan:  "",
				},
				"contextualTypeAlias": {
					text:           "Cb(): void",
					parameterCount: 0,
					parameterSpan:  "",
				},
				"contextualFunctionType": {
					text:           "cb2(): void",
					parameterCount: 0,
					parameterSpan:  "",
				},
			},
		},
		{
			title: "signatureHelpAnonymousFunction",
			input: `var anonymousFunctionTest = function(n: number, s: string): (a: number, b: string) => string {
	return null;
}
anonymousFunctionTest(5, "")(/*anonymousFunction1*/1, /*anonymousFunction2*/"");`,
			output: map[string]verifySignatureHelpOptions{
				"anonymousFunction1": {
					text:           `(a: number, b: string): string`,
					parameterCount: 2,
					parameterSpan:  "a: number",
				},
				"anonymousFunction2": {
					text:           `(a: number, b: string): string`,
					parameterCount: 2,
					parameterSpan:  "b: string",
				},
			},
		},
		{
			title: "signatureHelpAtEOFs",
			input: `function Foo(arg1: string, arg2: string) {
}

Foo(/**/`,
			output: map[string]verifySignatureHelpOptions{
				"": {
					text:           "Foo(arg1: string, arg2: string): void",
					parameterCount: 2,
					parameterSpan:  "arg1: string",
				},
			},
		},
		{
			title: "signatureHelpBeforeSemicolon1",
			input: `function Foo(arg1: string, arg2: string) {
}

Foo(/**/;`,
			output: map[string]verifySignatureHelpOptions{
				"": {
					text:           "Foo(arg1: string, arg2: string): void",
					parameterCount: 2,
					parameterSpan:  "arg1: string",
				},
			},
		},
		{
			title: "signatureHelpCallExpression",
			input: `function fnTest(str: string, num: number) { }
fnTest(/*1*/'', /*2*/5);`,
			output: map[string]verifySignatureHelpOptions{
				"1": {
					text:           `fnTest(str: string, num: number): void`,
					parameterCount: 2,
					parameterSpan:  "str: string",
				},
				"2": {
					text:           `fnTest(str: string, num: number): void`,
					parameterCount: 2,
					parameterSpan:  "num: number",
				},
			},
		},
		{
			title: "signatureHelpConstructExpression",
			input: `class sampleCls { constructor(str: string, num: number) { } }
var x = new sampleCls(/*1*/"", /*2*/5);`,
			output: map[string]verifySignatureHelpOptions{
				"1": {
					text:           "sampleCls(str: string, num: number): sampleCls",
					parameterCount: 2,
					parameterSpan:  "str: string",
				},
				"2": {
					text:           "sampleCls(str: string, num: number): sampleCls",
					parameterCount: 2,
					parameterSpan:  "num: number",
				},
			},
		},
		{
			title: "signatureHelpConstructorInheritance",
			input: `class base {
constructor(s: string);
constructor(n: number);
constructor(a: any) { }
}
class B1 extends base { }
class B2 extends B1 { }
class B3 extends B2 {
    constructor() {
        super(/*indirectSuperCall*/3);
    }
}`,
			output: map[string]verifySignatureHelpOptions{
				"indirectSuperCall": {
					text:           "B2(n: number): B2",
					parameterCount: 1,
					parameterSpan:  "n: number",
				},
			},
		},
		{
			title: "signatureHelpConstructorOverload",
			input: `class clsOverload { constructor(); constructor(test: string); constructor(test?: string) { } }
var x = new clsOverload(/*1*/);
var y = new clsOverload(/*2*/'');`,
			output: map[string]verifySignatureHelpOptions{
				"1": {
					text:           "clsOverload(): clsOverload",
					parameterCount: 0,
				},
				"2": {
					text:           "clsOverload(test: string): clsOverload",
					parameterCount: 1,
					parameterSpan:  "test: string",
				},
			},
		},
		{
			title: "signatureHelpEmptyLists",
			input: `function Foo(arg1: string, arg2: string) {
}

Foo(/*1*/);
function Bar<T>(arg1: string, arg2: string) { }
Bar</*2*/>();`,
			output: map[string]verifySignatureHelpOptions{
				"1": {
					text:           "Foo(arg1: string, arg2: string): void",
					parameterCount: 2,
					parameterSpan:  "arg1: string",
				},
				"2": {
					text:           "Bar<T>(arg1: string, arg2: string): void",
					parameterCount: -1,
					parameterSpan:  "T",
				},
			},
		},
		{
			title: "signatureHelpExpandedRestTuples",
			input: `export function complex(item: string, another: string, ...rest: [] | [settings: object, errorHandler: (err: Error) => void] | [errorHandler: (err: Error) => void, ...mixins: object[]]) {
    
}

complex(/*1*/);
complex("ok", "ok", /*2*/);
complex("ok", "ok", e => void e, {}, /*3*/);`,

			output: map[string]verifySignatureHelpOptions{
				"1": {
					text:           "complex(item: string, another: string): void",
					parameterCount: 2,
					parameterSpan:  "item: string",
					isVariadic:     false,
				},
				"2": {
					text:           "complex(item: string, another: string, settings: object, errorHandler: object): void", // Needs createSignatureHelpParameterForParameter
					parameterCount: 4,
					parameterSpan:  "settings: object",
					isVariadic:     false,
				},
				"3": {
					text:           "complex(item: string, another: string, errorHandler: object, ...mixins: object): void", // Needs createSignatureHelpParameterForParameter
					parameterCount: 4,
					parameterSpan:  "...mixins: object",
					isVariadic:     true,
				},
			},
		},
		{
			title: "signatureHelpExpandedRestUnlabeledTuples",
			input: `export function complex(item: string, another: string, ...rest: [] | [object, (err: Error) => void] | [(err: Error) => void, ...object[]]) {
   
}

complex(/*1*/);
complex("ok", "ok", /*2*/);
complex("ok", "ok", e => void e, {}, /*3*/);`,
			output: map[string]verifySignatureHelpOptions{
				"1": {
					text:           "complex(item: string, another: string): void",
					parameterCount: 2,
					parameterSpan:  "item: string",
					isVariadic:     false,
				},
				"2": {
					text:           "complex(item: string, another: string, rest_0: object, rest_1: object): void", // Needs createSignatureHelpParameterForParameter
					parameterCount: 4,
					parameterSpan:  "rest_0: object",
					isVariadic:     false,
				},
				"3": {
					text:           "complex(item: string, another: string, rest_0: object, ...rest: object): void", // Needs createSignatureHelpParameterForParameter
					parameterCount: 4,
					parameterSpan:  "...rest: object",
					isVariadic:     true,
				},
			},
		},
		{
			title: "signatureHelpExpandedTuplesArgumentIndex",
			input: `function foo(...args: [string, string] | [number, string, string]
) {

}

foo(""/*2*/, ""/*3*/)
foo(123/*4*/, ""/*5*/, )
foo(123/*6*/, ""/*7*/, ""/*8*/)`,
			output: map[string]verifySignatureHelpOptions{
				"1": { // problem with foo(123/*1*/,)
					text:                      "foo(args_0: number, args_1: string, args_2: string): void",
					parameterCount:            3,
					parameterSpan:             "args_0: number",
					isVariadic:                false,
					overrideSelectedItemIndex: 1,
				},
				"2": {
					text:                      "foo(args_0: string, args_1: string): void",
					parameterCount:            2,
					parameterSpan:             "args_0: string",
					isVariadic:                false,
					overrideSelectedItemIndex: 0,
				},
				"3": {
					text:                      "foo(args_0: string, args_1: string): void",
					parameterCount:            2,
					parameterSpan:             "args_1: string",
					isVariadic:                false,
					overrideSelectedItemIndex: 0,
				},
				"4": {
					text: "foo(args_0: number, args_1: string, args_2: string): void",

					parameterCount:            3,
					parameterSpan:             "args_0: number",
					isVariadic:                false,
					overrideSelectedItemIndex: 1,
				},
				"5": {
					text:                      "foo(args_0: number, args_1: string, args_2: string): void",
					parameterCount:            3,
					parameterSpan:             "args_1: string",
					isVariadic:                false,
					overrideSelectedItemIndex: 1,
				},
				"6": {
					text:                      "foo(args_0: number, args_1: string, args_2: string): void",
					parameterCount:            3,
					parameterSpan:             "args_0: number",
					isVariadic:                false,
					overrideSelectedItemIndex: 1,
				},
				"7": {
					text:                      "foo(args_0: number, args_1: string, args_2: string): void",
					parameterCount:            3,
					parameterSpan:             "args_1: string",
					isVariadic:                false,
					overrideSelectedItemIndex: 1,
				},
				"8": {
					text:                      "foo(args_0: number, args_1: string, args_2: string): void",
					parameterCount:            3,
					parameterSpan:             "args_2: string",
					isVariadic:                false,
					overrideSelectedItemIndex: 1,
				},
			},
		},
		{
			title: "signatureHelpExplicitTypeArguments",
			input: `declare function f<T = boolean, U = string>(x: T, y: U): T;
f<number, string>(/*1*/);
f(/*2*/);
f<number>(/*3*/);
f<number, string, boolean>(/*4*/);

interface A { a: number }
interface B extends A { b: string }
declare function g<T, U, V extends A = B>(x: T, y: U, z: V): T;
declare function h<T, U, V extends A>(x: T, y: U, z: V): T;
declare function j<T, U, V = B>(x: T, y: U, z: V): T;
g(/*5*/);
h(/*6*/);
j(/*7*/);
g<number>(/*8*/);
h<number>(/*9*/);
j<number>(/*10*/);`,
			output: map[string]verifySignatureHelpOptions{
				"1": {text: "f(x: number, y: string): number", parameterCount: 2, parameterSpan: "x: number"},
				"2": {text: "f(x: boolean, y: string): boolean", parameterCount: 2, parameterSpan: "x: boolean"},
				// too few -- fill in rest with default
				"3": {text: "f(x: number, y: string): number", parameterCount: 2, parameterSpan: "x: number"},
				// too many -- ignore extra type arguments
				"4": {text: "f(x: number, y: string): number", parameterCount: 2, parameterSpan: "x: number"},

				// not matched signature and no type arguments
				"5": {text: "g(x: unknown, y: unknown, z: object): unknown", parameterCount: 3, parameterSpan: "x: unknown"},
				"6": {text: "h(x: unknown, y: unknown, z: object): unknown", parameterCount: 3, parameterSpan: "x: unknown"},
				"7": {text: "j(x: unknown, y: unknown, z: object): unknown", parameterCount: 3, parameterSpan: "x: unknown"},
				// not matched signature and too few type arguments
				"8":  {text: "g(x: number, y: unknown, z: object): number", parameterCount: 3, parameterSpan: "x: number"},
				"9":  {text: "h(x: number, y: unknown, z: object): number", parameterCount: 3, parameterSpan: "x: number"},
				"10": {text: "j(x: number, y: unknown, z: object): number", parameterCount: 3, parameterSpan: "x: number"},
			},
		},
	}

	for _, rec := range testCases {
		t.Run(rec.title, func(t *testing.T) {
			t.Parallel()
			testData := parseTestdata("/file1.ts", rec.input, "/file1.ts")
			service := createLanguageService(testData.files[0].filename, map[string]string{
				testData.files[0].filename: testData.files[0].content,
			})
			file := parser.ParseSourceFile(testData.files[0].filename, tspath.Path(testData.files[0].filename), testData.files[0].content, core.ScriptTargetLatest, scanner.JSDocParsingModeParseAll)

			markerNumber := 0
			for i, marker := range testData.markerPositions {
				result := service.GetSignatureHelpItems(file.FileName(), marker.position, nil)
				if result == nil {
					t.Fatal("expected result to be non-nil")
				}
				if _, exists := testData.markerPositions[i]; !exists {
					t.Fatal("marker not found in test data")
				}
				assert.Equal(t, rec.output[i].text, result.Signatures[result.ActiveSignature].Label)
				if rec.output[i].parameterCount != -1 {
					assert.Equal(t, rec.output[i].parameterCount, len(*result.Signatures[result.ActiveSignature].Parameters))
				}

				if len(*result.Signatures[result.ActiveSignature].Parameters) <= result.ActiveParameter || len(*result.Signatures[result.ActiveSignature].Parameters) == 0 {
					assert.Equal(t, rec.output[i].parameterSpan, "")
				} else {
					assert.Equal(t, rec.output[i].parameterSpan, (*result.Signatures[result.ActiveSignature].Parameters)[result.ActiveParameter].Label)
				}
				markerNumber++
			}
		})
	}
}

func createLanguageService(fileName string, files map[string]string) *ls.LanguageService {
	projectService, _ := projecttestutil.Setup(files)
	projectService.OpenFile(fileName, files[fileName], core.ScriptKindTS, "")
	project := projectService.Projects()[0]
	return project.LanguageService()
}
