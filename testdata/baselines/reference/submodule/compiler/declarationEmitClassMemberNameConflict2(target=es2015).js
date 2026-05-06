//// [tests/cases/compiler/declarationEmitClassMemberNameConflict2.ts] ////

//// [declarationEmitClassMemberNameConflict2.ts]
const Bar = 'bar';

enum Hello {
    World
}

enum Hello1 {
    World1
}

class Foo {
    // Same names + string => OK
    Bar = Bar;

    // Same names + enum => OK
    Hello = Hello;

    // Different names + enum => OK
    Hello2 = Hello1;
}

//// [declarationEmitClassMemberNameConflict2.js]
"use strict";
const Bar = 'bar';
var Hello;
(function (Hello) {
    Hello[Hello["World"] = 0] = "World";
})(Hello || (Hello = {}));
var Hello1;
(function (Hello1) {
    Hello1[Hello1["World1"] = 0] = "World1";
})(Hello1 || (Hello1 = {}));
class Foo {
    constructor() {
        // Same names + string => OK
        this.Bar = Bar;
        // Same names + enum => OK
        this.Hello = Hello;
        // Different names + enum => OK
        this.Hello2 = Hello1;
    }
}


//// [declarationEmitClassMemberNameConflict2.d.ts]
const Bar = "bar";
enum Hello {
    World = 0
}
enum Hello1 {
    World1 = 0
}
class Foo {
    Bar: string;
    Hello: typeof Hello;
    Hello2: typeof Hello1;
}


//// [DtsFileErrors]


declarationEmitClassMemberNameConflict2.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declarationEmitClassMemberNameConflict2.d.ts (1 errors) ====
    const Bar = "bar";
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    enum Hello {
        World = 0
    }
    enum Hello1 {
        World1 = 0
    }
    class Foo {
        Bar: string;
        Hello: typeof Hello;
        Hello2: typeof Hello1;
    }
    