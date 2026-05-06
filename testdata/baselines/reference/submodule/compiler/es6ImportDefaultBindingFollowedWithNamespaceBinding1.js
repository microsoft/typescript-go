//// [tests/cases/compiler/es6ImportDefaultBindingFollowedWithNamespaceBinding1.ts] ////

//// [es6ImportDefaultBindingFollowedWithNamespaceBinding_0.ts]
var a = 10;
export default a;

//// [es6ImportDefaultBindingFollowedWithNamespaceBinding_1.ts]
import defaultBinding, * as nameSpaceBinding  from "./es6ImportDefaultBindingFollowedWithNamespaceBinding_0";
var x: number = defaultBinding;

//// [es6ImportDefaultBindingFollowedWithNamespaceBinding_0.js]
var a = 10;
export default a;
//// [es6ImportDefaultBindingFollowedWithNamespaceBinding_1.js]
import defaultBinding from "./es6ImportDefaultBindingFollowedWithNamespaceBinding_0";
var x = defaultBinding;


//// [es6ImportDefaultBindingFollowedWithNamespaceBinding_0.d.ts]
var a: number;
export default a;
//// [es6ImportDefaultBindingFollowedWithNamespaceBinding_1.d.ts]
export {};


//// [DtsFileErrors]


es6ImportDefaultBindingFollowedWithNamespaceBinding_0.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== es6ImportDefaultBindingFollowedWithNamespaceBinding_0.d.ts (1 errors) ====
    var a: number;
    ~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    export default a;
    
==== es6ImportDefaultBindingFollowedWithNamespaceBinding_1.d.ts (0 errors) ====
    export {};
    