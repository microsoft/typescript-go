package ls_test

import (
	"sort"
	"testing"

	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/ls"
	"github.com/microsoft/typescript-go/internal/testutil/projecttestutil"
	"gotest.tools/v3/assert"
)

func runFoldingRangeTest(t *testing.T, input string) {
	testData := fourslash.ParseTestData(t, input, "/file1.ts")
	markerPositions := testData.Ranges
	ctx := projecttestutil.WithRequestID(t.Context())
	service, done := createLanguageService(ctx, testData.Files[0].FileName(), map[string]string{
		testData.Files[0].FileName(): testData.Files[0].Content,
	})
	defer done()

	foldingRanges := service.ProvideFoldingRange(ctx, ls.FileNameToDocumentURI("/file1.ts"))
	if len(foldingRanges) != len(markerPositions) {
		t.Fatalf("Expected %d folding ranges, got %d", len(markerPositions), len(foldingRanges))
	}
	sort.Slice(markerPositions, func(i, j int) bool {
		if markerPositions[i].LSRange.Start.Line != markerPositions[j].LSRange.Start.Line {
			return markerPositions[i].LSRange.Start.Line < markerPositions[j].LSRange.Start.Line
		}
		if markerPositions[i].LSRange.End.Line != markerPositions[j].LSRange.End.Line {
			return markerPositions[i].LSRange.End.Line < markerPositions[j].LSRange.End.Line
		} else if markerPositions[i].LSRange.Start.Character != markerPositions[j].LSRange.Start.Character {
			return markerPositions[i].LSRange.Start.Character < markerPositions[j].LSRange.Start.Character
		}
		return markerPositions[i].LSRange.End.Character < markerPositions[j].LSRange.End.Character
	})
	for i, marker := range markerPositions {
		assert.DeepEqual(t, marker.LSRange.Start.Line, foldingRanges[i].StartLine)
		assert.DeepEqual(t, marker.LSRange.End.Line, foldingRanges[i].EndLine)
		assert.DeepEqual(t, marker.LSRange.Start.Character, *foldingRanges[i].StartCharacter)
		assert.DeepEqual(t, marker.LSRange.End.Character, *foldingRanges[i].EndCharacter)
	}
}

func TestFolding(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		title             string
		input             string
		expectedLocations map[string]*collections.Set[string]
	}{
		{
			title: "getOutliningSpansForRegionsNoSingleLineFolds",
			input: `[|//#region
function foo()[| {

}|]
[|//these
//should|]
//#endregion not you|]
[|// be
// together|]

[|//#region bla bla bla

function bar()[| { }|]

//#endregion|]`,
		},
		{
			title: "getOutliningSpansForComments",
			input: `[|/*
   Block comment at the beginning of the file before module:
       line one of the comment
       line two of the comment
       line three
       line four
       line five
*/|]
declare module "m";
[|// Single line comments at the start of the file
// line 2
// line 3
// line 4|]
declare module "n";`,
		},
		{
			title: "getOutliningSpansForRegions",
			input: `// region without label
		[|// #region

		// #endregion|]

		// region without label with trailing spaces
		[|// #region

		// #endregion|]

		// region with label
		[|// #region label1

		// #endregion|]

		// region with extra whitespace in all valid locations
		            [|//              #region          label2    label3

		       //        #endregion|]

		// No space before directive
		[|//#region label4

		//#endregion|]

		// Nested regions
		[|// #region outer

		[|// #region inner

		// #endregion inner|]

		// #endregion outer|]

		// region delimiters not valid when there is preceding text on line
		test // #region invalid1

		test // #endregion`,
		},
		{
			title: "outliningSpansSwitchCases",
			input: `switch (undefined)[| {
case 0:[|
  console.log(1)
  console.log(2)
  break;
  console.log(3);|]
case 1:[|
  break;|]
case 2:[|
  break;
  console.log(3);|]
case 3:[|
  console.log(4);|]

case 4:
case 5:
case 6:[|


  console.log(5);|]

case 7:[| console.log(6);|]

case 8:[| [|{
  console.log(8);
  break;
}|]
console.log(8);|]

default:[|
  console.log(7);
  console.log(8);|]
}|]`,
		},
		{
			title: "outliningSpansForParenthesizedExpression",
			input: `const a = [|(
   true
       ? true
       : false
           ? true
           : false
)|];

const b = ( 1 );

const c = [|(
   1
)|];

( 1 );

[|(
   [|(
       [|(
           1
       )|]
   )|]
)|];

[|(
   [|(
       ( 1 )
   )|]
)|];`,
		},
		{
			title: "outliningSpansForInportsAndExports",
			input: `import { a1, a2 } from "a";
;
import {
} from "a";
;
import [|{
  b1,
  b2,
}|] from "b";
;
import j1 from "./j" assert { type: "json" };
;
import j2 from "./j" assert {
};
;
import j3 from "./j" assert [|{
  type: "json"
}|];
;
[|import { a5, a6 } from "a";
import [|{
  a7,
  a8,
}|] from "a";|]
export { a1, a2 };
;
export { a3, a4 } from "a";
;
export {
};
;
export [|{
  b1,
  b2,
}|];
;
export {
} from "b";
;
export [|{
  b3,
  b4,
}|] from "b";
;`,
		},
		{
			title: "outliningSpansForImportAndExportAttributes",
			input: `import { a1, a2 } from "a";
;
import {
} from "a";
;
import [|{
  b1,
  b2,
}|] from "b";
;
import j1 from "./j" with { type: "json" };
;
import j2 from "./j" with {
};
;
import j3 from "./j" with [|{
  type: "json"
}|];
;
[|import { a5, a6 } from "a";
import [|{
  a7,
  a8,
}|] from "a";|]
export { a1, a2 };
;
export { a3, a4 } from "a";
;
export {
};
;
export [|{
  b1,
  b2,
}|];
;
export {
} from "b";
;
export [|{
  b3,
  b4,
}|] from "b";
;`,
		},
		{
			title: "outliningSpansForFunction",
			input: `[|(
   a: number,
   b: number
) => {
   return a + b;
}|];

(a: number, b: number) =>[| {
   return a + b;
}|]

const f1 = function[| (
   a: number
   b: number
) {
   return a + b;
}|]

const f2 = function (a: number, b: number)[| {
   return a + b;
}|]

function f3[| (
   a: number
   b: number
) {
   return a + b;
}|]

function f4(a: number, b: number)[| {
   return a + b;
}|]

class Foo[| {
   constructor[|(
       a: number,
       b: number
   ) {
       this.a = a;
       this.b = b;
   }|]

   m1[|(
       a: number,
       b: number
   ) {
       return a + b;
   }|]

   m1(a: number, b: number)[| {
       return a + b;
   }|]
}|]

declare function foo(props: any): void;
foo[|(
   a =>[| {

   }|]
)|]

foo[|(
   (a) =>[| {

   }|]
)|]

foo[|(
   (a, b, c) =>[| {

   }|]
)|]

foo[|([|
   (a,
    b,
    c) => {

   }|]
)|]`,
		},
		{
			title: "outliningSpansForArrowFunctionBody",
			input: `() => 42;
() => ( 42 );
() =>[| {
    42
}|];
() => [|(
    42
)|];
() =>[| "foo" +
    "bar" +
    "baz"|];`,
		},
		{
			title: "outliningSpansForArguments",
			input: `console.log(123, 456);
console.log(
);
console.log[|(
    123, 456
)|];
console.log[|(
    123,
    456
)|];
() =>[| console.log[|(
    123,
    456
)|]|];`,
		},
		{
			title: "outliningForNonCompleteInterfaceDeclaration",
			input: `interface I`,
		},
		{
			title: "incrementalParsingWithJsDoc",
			input: `[|import a from 'a/aaaaaaa/aaaaaaa/aaaaaa/aaaaaaa';
import b from 'b';
import c from 'c';|]

[|/** @internal */|]
export class LanguageIdentifier[| { }|]`,
		},
		{
			title: "incrementalParsingWithJsDoc_2",
			input: `[|import a from 'a/aaaaaaa/aaaaaaa/aaaaaa/aaaaaaa';
/**/import b from 'b';
import c from 'c';|]

[|/** @internal */|]
export class LanguageIdentifier[| { }|]`,
		},
		{
			title: "getOutliningSpansForUnbalancedRegion",
			input: `// top-heavy region balance
// #region unmatched

[|// #region matched

// #endregion matched|]`,
		},
		{
			title: "getOutliningSpansForTemplateLiteral",
			input: "declare function tag(...args: any[]): void\nconst a = [|`signal line`|]\nconst b = [|`multi\nline`|]\nconst c = tag[|`signal line`|]\nconst d = tag[|`multi\nline`|]\nconst e = [|`signal ${1} line`|]\nconst f = [|`multi\n${1}\nline`|]\nconst g = tag[|`signal ${1} line`|]\nconst h = tag[|`multi\n${1}\nline`|]\nconst i = ``",
		},
		{
			title: "getOutliningSpansForImports",
			input: `[|import * as ns from "mod";

import d from "mod";
import { a, b, c } from "mod";

import r = require("mod");|]

// statement
var x = 0;

// another set of imports
[|import * as ns from "mod";
import d from "mod";
import { a, b, c } from "mod";
import r = require("mod");|]`,
		},
		{
			title: "getOutliningSpansDepthElseIf",
			input: `if (1)[| {
   1;
}|] else if (1)[| {
   1;
}|] else if (1)[| {
   1;
}|] else if (1)[| {
   1;
}|] else if (1)[| {
   1;
}|] else if (1)[| {
   1;
}|] else if (1)[| {
   1;
}|] else if (1)[| {
   1;
}|] else if (1)[| {
   1;
}|] else if (1)[| {
   1;
}|] else if (1)[| {
   1;
}|] else if (1)[| {
   1;
}|] else if (1)[| {
   1;
}|] else if (1)[| {
   1;
}|] else if (1)[| {
   1;
}|] else if (1)[| {
   1;
}|] else if (1)[| {
   1;
}|] else if (1)[| {
   1;
}|] else if (1)[| {
   1;
}|] else if (1)[| {
   1;
}|] else if (1)[| {
   1;
}|] else if (1)[| {
   1;
}|] else if (1)[| {
   1;
}|] else if (1)[| {
   1;
}|] else if (1)[| {
   1;
}|] else if (1)[| {
   1;
}|] else if (1)[| {
   1;
}|] else if (1)[| {
   1;
}|] else if (1)[| {
   1;
}|] else if (1)[| {
   1;
}|] else if (1)[| {
   1;
}|] else if (1)[| {
   1;
}|] else if (1)[| {
   1;
}|] else if (1)[| {
   1;
}|] else if (1)[| {
   1;
}|] else if (1)[| {
   1;
}|] else if (1)[| {
   1;
}|] else if (1)[| {
   1;
}|] else if (1)[| {
   1;
}|] else if (1)[| {
   1;
}|] else[| {
   1;
}|]`,
		},
		{
			title: "getOutliningSpansDepthChainedCalls",
			input: `declare var router: any;
router
   .get[|("/", async(ctx) =>[|{
       ctx.body = "base";
   }|])|]
   .post[|("/a", async(ctx) =>[|{
       //a
   }|])|]
   .get[|("/", async(ctx) =>[|{
       ctx.body = "base";
   }|])|]
   .post[|("/a", async(ctx) =>[|{
       //a
   }|])|]
   .get[|("/", async(ctx) =>[|{
       ctx.body = "base";
   }|])|]
   .post[|("/a", async(ctx) =>[|{
       //a
   }|])|]
   .get[|("/", async(ctx) =>[|{
       ctx.body = "base";
   }|])|]
   .post[|("/a", async(ctx) =>[|{
       //a
   }|])|]
   .get[|("/", async(ctx) =>[|{
       ctx.body = "base";
   }|])|]
   .post[|("/a", async(ctx) =>[|{
       //a
   }|])|]
   .get[|("/", async(ctx) =>[|{
       ctx.body = "base";
   }|])|]
   .post[|("/a", async(ctx) =>[|{
       //a
   }|])|]
   .get[|("/", async(ctx) =>[|{
       ctx.body = "base";
   }|])|]
   .post[|("/a", async(ctx) =>[|{
       //a
   }|])|]
   .get[|("/", async(ctx) =>[|{
       ctx.body = "base";
   }|])|]
   .post[|("/a", async(ctx) =>[|{
       //a
   }|])|]
   .get[|("/", async(ctx) =>[|{
       ctx.body = "base";
   }|])|]
   .post[|("/a", async(ctx) =>[|{
       //a
   }|])|]
   .get[|("/", async(ctx) =>[|{
       ctx.body = "base";
   }|])|]
   .post[|("/a", async(ctx) =>[|{
       //a
   }|])|]
   .get[|("/", async(ctx) =>[|{
       ctx.body = "base";
   }|])|]
   .post[|("/a", async(ctx) =>[|{
       //a
   }|])|]
   .get[|("/", async(ctx) =>[|{
       ctx.body = "base";
   }|])|]
   .post[|("/a", async(ctx) =>[|{
       //a
   }|])|]
   .get[|("/", async(ctx) =>[|{
       ctx.body = "base";
   }|])|]
   .post[|("/a", async(ctx) =>[|{
       //a
   }|])|]
   .get[|("/", async(ctx) =>[|{
       ctx.body = "base";
   }|])|]
   .post[|("/a", async(ctx) =>[|{
       //a
   }|])|]
   .get[|("/", async(ctx) =>[|{
       ctx.body = "base";
   }|])|]
   .post[|("/a", async(ctx) =>[|{
       //a
   }|])|]
   .get[|("/", async(ctx) =>[|{
       ctx.body = "base";
   }|])|]
   .post[|("/a", async(ctx) =>[|{
       //a
   }|])|]
   .get[|("/", async(ctx) =>[|{
       ctx.body = "base";
   }|])|]
   .post[|("/a", async(ctx) =>[|{
       //a
   }|])|]
   .get[|("/", async(ctx) =>[|{
       ctx.body = "base";
   }|])|]
   .post[|("/a", async(ctx) =>[|{
       //a
   }|])|]`,
		},
		{
			title: "getOutliningSpans",
			input: `// interface
interface IFoo[| {
   getDist(): number;
}|]

// class members
class Foo[| {
   constructor()[| {
   }|]

   public foo(): number[| {
       return 0;
   }|]

   public get X()[| {
       return 1;
   }|]

   public set X(v: number)[| {
   }|]

   public member = function f()[| {

   }|]
}|]
// class expressions
[|(new class[| {
    bla()[| {

    }|]
}|])|]
switch(1)[| {
case 1:[| break;|]
}|]

var array =[| [
   1,
   2
]|]

// modules
module m1[| {
   module m2[| { }|]
   module m3[| {
       function foo()[| {

       }|]

       interface IFoo2[| {

       }|]

       class foo2 implements IFoo2[| {

       }|]
   }|]
}|]

// function declaration
function foo(): number[| {
   return 0;
}|]

// function expressions
[|(function f()[| {

}|])|]

// trivia handeling
class ClassFooWithTrivia[| /*  some comments */
  /* more trivia */ {


   [|/*some trailing trivia */|]
}|] /* even more */

// object literals
var x =[|{
 a:1,
 b:2,
 get foo()[| {
   return 1;
 }|]
}|]
//outline with deep nesting
var nest =[| [[|[[|[[|[[|[[|[[|[[|[[|[[|[[|[[|[[|[[|[[|[[|[[|[[|[[|[[|[[|[[|[[|[[|[[|[[|[[|[[|[[|[[|[[|[[|[[|[[|[[|[[|[
 [|[
     [
       [
         [
           [
             1,2,3
           ]
         ]
       ]
     ]
 ]|]
]|]]|]]|]]|]]|]]|]]|]]|]]|]]|]]|]]|]]|]]|]]|]]|]]|]]|]]|]]|]]|]]|]]|]]|]]|]]|]]|]]|]]|]]|]]|]]|]]|]]|]]|]]|];

//outline after a deeply nested node
class AfterNestedNodes[| {
}|]
// function arguments
function f(x: number[], y: number[])[| {
   return 3;
}|]
f[|(
//  single line array literal span won't render in VS
   [|[0]|],
   [|[
       1,
       2
   ]|]
)|];

class C<T>[| {
   foo: T;
}|]

class D<T> extends C<T>[| {
   constructor(x)[| {
       super<T>(x);
   }|]
}|]`,
		},
		{
			title: "getOutliningForTypeLiteral",
			input: `type A =[| {
   a: number;
}|]

type B =[| {
  a:[| {
      a1:[| {
          a2:[| {
              x: number;
              y: number;
          }|]
      }|]
  }|],
  b:[| {
      x: number;
  }|],
  c:[| {
      x: number;
  }|]
}|]`,
		},
		{
			title: "getOutliningForTupleType",
			input: `type A =[| [
   number,
   number,
   number
]|]

type B =[| [
   [|[
       [|[
           number,
           number,
           number
       ]|]
   ]|]
]|]`,
		},
		{
			title: "getOutliningForSingleLineComments",
			input: `[|// Single line comments at the start of the file
		// line 2
		// line 3
		// line 4|]
		module Sayings[| {

		   [|/*
		   */|]
		   [|// A sequence of
		   // single line|]
		   [|/*
		       and block
		   */|]
		   [|// comments
		   //|]
		   export class Sample[| {
		   }|]
		}|]

		interface IFoo[| {
		   [|// all consecutive single line comments should be in one block regardless of their number or empty lines/spaces in between

		   // comment 2
		   // comment 3

		   //comment 4
		   /// comment 5
		   ///// comment 6

		   //comment 7
		   ///comment 8
		   // comment 9
		   // //comment 10

		   // // //comment 11
		   // comment 12
		   // comment 13
		   // comment 14
		   // comment 15

		   // comment 16
		   // comment 17
		   // comment 18
		   // comment 19
		   // comment 20
		   // comment 21|]

		   getDist(): number; // One single line comment should not be collapsed
		}|]

		// One single line comment should not be collapsed
		class WithOneSingleLineComment[| {
		}|]

		function Foo()[| {
		  [|// comment 1
		    // comment 2|]
		   this.method = function (param)[| {
		   }|]

		  [|// comment 1
		    // comment 2|]
		   function method(param)[| {
		   }|]
		}|]`,
		},
		{
			title: "getOutliningForObjectsInArray",
			input: `// objects in x should generate outlining spans that do not render in VS
const x =[| [
    [|{ a: 0 }|],
    [|{ b: 1 }|],
    [|{ c: 2 }|]
]|];
// objects in y should generate outlining spans that render as expected
const y =[| [
    [|{
        a: 0
    }|],
    [|{
        b: 1
    }|],
    [|{
        c: 2
    }|]
]|];
// same behavior for nested arrays
const w =[| [
    [|[ 0 ]|],
    [|[ 1 ]|],
    [|[ 2 ]|]
]|];

const z =[| [
    [|[
        0
    ]|],
    [|[
        1
    ]|],
    [|[
        2
    ]|]
]|];
// multiple levels of nesting work as expected
const z =[| [
    [|[
        [|{ hello: 0 }|]
    ]|],
    [|[
        [|{ hello: 3 }|]
    ]|],
    [|[
        [|{ hello: 5 }|],
        [|{ hello: 7 }|]
    ]|]
]|];`,
		},
		{
			title: "getOutliningForObjectDestructuring",
			input: `const[| {
   a,
   b,
   c
}|] =[| {
   a: 1,
   b: 2,
   c: 3
}|]

const[| {
   a:[| {
       a_1,
       a_2,
       a_3:[| {
           a_3_1,
           a_3_2,
           a_3_3,
       }|],
   }|],
   b,
   c
}|] =[| {
   a:[| {
       a_1: 1,
       a_2: 2,
       a_3:[| {
           a_3_1: 1,
           a_3_2: 1,
           a_3_3: 1
       }|],
   }|],
   b: 2,
   c: 3
}|]`,
		},
		{
			title: "getOutliningForBlockComments",
			input: `[|/*
		   Block comment at the beginning of the file before module:
		       line one of the comment
		       line two of the comment
		       line three
		       line four
		       line five
		*/|]
		module Sayings[| {
		   [|/*
		   Comment before class:
		       line one of the comment
		       line two of the comment
		       line three
		       line four
		       line five
		   */|]
		   export class Greeter[| {
		       [|/*
		           Comment before a string identifier
		           line two of the comment
		       */|]
		       greeting: string;
		       [|/*
		           constructor
		           parameter message as a string
		       */|]

		       [|/*
		           Multiple comments should be collapsed individually
		       */|]
		       constructor(message: string /* do not collapse this */)[| {
		           this.greeting = message;
		       }|]
		       [|/*
		           method of a class
		       */|]
		       greet()[| {
		           return "Hello, " + this.greeting;
		       }|]
		   }|]
		}|]

		[|/*
		   Block comment for interface. The ending can be on the same line as the declaration.
		*/|]interface IFoo[| {
		   [|/*
		   Multiple block comments
		   */|]

		   [|/*
		   should be collapsed
		   */|]

		   [|/*
		   individually
		   */|]

		                                                                                                                             [|/*
		                                                                   this comment has trailing space before /* and after *-/ signs
		   */|]

		   [|/**
		    *
		    *
		    *
		    */|]

		   [|/*
		   */|]

		   [|/*
		   */|]
		   // single line comments in the middle should not have an effect
		   [|/*
		   */|]

		   [|/*
		   */|]

		   [|/*
		   this block comment ends
		   on the same line */|]  [|/* where the following comment starts
		       should be collapsed separately
		   */|]

		   getDist(): number;
		}|]

		var x =[|{
		 a:1,
		 b: 2,
		 [|/*
		       Over a function in an object literal
		 */|]
		 get foo()[| {
		   return 1;
		 }|]
		}|]

		// Over a function expression assigned to a variable
		[|/**
		 * Return a sum
		 * @param {Number} y
		 * @param {Number} z
		 * @returns {Number} the sum of y and z
		 */|]
		const sum2 = (y, z) =>[| {
		    return y + z;
		}|];

		// Over a variable
		[|/**
		* foo
		*/|]
		const foo = null;

		function Foo()[| {
		  [|/**
		    * Description
		    *
		    * @param {string} param
		    * @returns
		    */|]
		   this.method = function (param)[| {
		   }|]

		  [|/**
		    * Description
		    *
		    * @param {string} param
		    * @returns
		    */|]
		   function method(param)[| {
		   }|]
		}|]

		function fn1()[| {
		   [|/**
		    * comment
		    */|]
		}|]
		function fn2()[| {
		   [|/**
		    * comment
		    */|]

		   [|/**
		    * comment
		    */|]
		}|]
		function fn3()[| {
		   const x = 1;

		   [|/**
		    * comment
		    */|]

		   [|/**
		    * comment
		    */|]
		}|]
		function fn4()[| {
		   [|/**
		    * comment
		    */|]
		    const x = 1;

		   [|/**
		    * comment
		    */|]
		}|]
		function fn5()[| {
		   [|/**
		    * comment
		    */|]

		   [|/**
		    * comment
		    */|]
		    return 1;
		}|]
		function fn6()[| {
		   [|/**
		    * comment
		    */|]

		   [|/**
		    * comment
		    */|]
		   const x = 1;
		}|]

		[|/*
		comment
		*/|]

		f6();

		class C1[| {
		   [|/**
		    * comment
		    */|]

		   [|/**
		    * comment
		    */|]
		}|]
		class C2[| {
		   private prop = 1;
		   [|/**
		    * comment
		    */|]

		   [|/**
		    * comment
		    */|]
		}|]
		class C3[| {
		   [|/**
		    * comment
		    */|]

		   private prop = 1;
		   [|/**
		    * comment
		    */|]
		}|]
		class C4[| {
		   [|/**
		    * comment
		    */|]

		   [|/**
		    * comment
		    */|]
		   private prop = 1;
		}|]

		[|/*
		comment
		*/|]
		new C4();

		module M1[| {
		   [|/**
		    * comment
		    */|]

		   [|/**
		    * comment
		    */|]
		}|]
		module M2[| {
		   export const a = 1;
		   [|/**
		    * comment
		    */|]

		   [|/**
		    * comment
		    */|]
		}|]
		module M3[| {
		   [|/**
		    * comment
		    */|]
		   export const a = 1;

		   [|/**
		    * comment
		    */|]
		}|]
		module M4[| {
		   [|/**
		    * comment
		    */|]

		   [|/**
		    * comment
		    */|]
		   export const a = 1;
		}|]
		interface I1[| {
		   [|/**
		    * comment
		    */|]

		   [|/**
		    * comment
		    */|]
		}|]
		interface I2[| {
		   x: number;
		   [|/**
		    * comment
		    */|]

		   [|/**
		    * comment
		    */|]
		}|]
		interface I3[| {
		   [|/**
		    * comment
		    */|]
		   x: number;

		   [|/**
		    * comment
		    */|]
		}|]
		interface I4[| {
		   [|/**
		    * comment
		    */|]

		   [|/**
		    * comment
		    */|]
		   x: number;
		}|]
		[|{
		   [|/**
		    * comment
		    */|]

		   [|/**
		    * comment
		    */|]
		}|]`,
		},
		{
			title: "getOutliningForArrayDestructuring",
			input: `const[| [
   a,
   b,
   c
]|] =[| [
   1,
   2,
   3
]|];

const[| [
   [|[
       [|[
           [|[
               a,
               b,
               c
           ]|]
       ]|]
   ]|],
   [|[
       a1,
       b1,
       c1
   ]|]
]|] =[| [
   [|[
       [|[
           [|[
               1,
               2,
               3
           ]|]
       ]|]
   ]|],
   [|[
       1,
       2,
       3
   ]|]
]|]`,
		},
		// {
		// 	title: "getJSXOutliningSpans",
		// 	input: `import React, { Component } from 'react';

		// export class Home extends Component[| {
		//  render()[| {
		//    return [|(
		//    [|<div>
		//      [|<h1>Hello, world!</h1>|]
		//      [|<ul>
		//        [|<li>
		//          [|<a [|href='https://get.asp.net/'|]>
		//            ASP.NET Core
		//          </a>|]
		//        </li>|]
		//        [|<li>[|<a [|href='https://facebook.github.io/react/'|]>React</a>|] for client-side code</li>|]
		//        [|<li>[|<a [|href='http://getbootstrap.com/'|]>Bootstrap</a>|] for layout and styling</li>|]
		//      </ul>|]
		//      <div
		//        [|accesskey="test"
		//        class="active"
		//        dir="auto"|] />
		//      <PageHeader [|title="Log in"
		//        {...[|{
		//          item: true,
		//          xs: 9,
		//          md: 5
		//        }|]}|]
		//      />
		//      [|<>
		//          text
		//      </>|]
		//    </div>|]
		//    )|];
		//  }|]
		// }|]`,
		// },
		{
			title: "corruptedTryExpressionsDontCrashGettingOutlineSpans",
			input: `try[| {
  var x = [
    {% try %}|]{% except %} 
  ]
} catch (e)[| {
  
}|]`,
		},
		{
			title: "outliningSpansForFunctions",
			input: `namespace NS[| {
   function f(x: number, y: number)[| {
       return x + y;
   }|]

   function g[|(
       x: number,
       y: number,
   ): number {
       return x + y;
   }|]
}|]`,
		},
		{
			title: "outliningSpansTrailingBlockCmmentsAfterStatements",
			input: `console.log(0);
[|/*
/ * Some text
 */|]`,
		},
		{
			title: "outlineSpansBlockCommentsWithoutStatements",
			input: `[|/*
/ * Some text
 */|]`,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.title, func(t *testing.T) {
			t.Parallel()
			runFoldingRangeTest(t, testCase.input)
		})
	}
}
