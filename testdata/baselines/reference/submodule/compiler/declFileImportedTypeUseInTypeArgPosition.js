//// [tests/cases/compiler/declFileImportedTypeUseInTypeArgPosition.ts] ////

//// [declFileImportedTypeUseInTypeArgPosition.ts]
class List<T> { }
declare module 'mod1' {
    class Foo {
    }
}

declare module 'moo' {
    import x = require('mod1');
    export var p: List<x.Foo>;
}




//// [declFileImportedTypeUseInTypeArgPosition.js]
"use strict";
class List {
}


//// [declFileImportedTypeUseInTypeArgPosition.d.ts]
class List<T> {
}
module 'mod1' {
    class Foo {
    }
}
module 'moo' {
    import x = require('mod1');
    var p: List<x.Foo>;
}


//// [DtsFileErrors]


declFileImportedTypeUseInTypeArgPosition.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declFileImportedTypeUseInTypeArgPosition.d.ts (1 errors) ====
    class List<T> {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    }
    module 'mod1' {
        class Foo {
        }
    }
    module 'moo' {
        import x = require('mod1');
        var p: List<x.Foo>;
    }
    