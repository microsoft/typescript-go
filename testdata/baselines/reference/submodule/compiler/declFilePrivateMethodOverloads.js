//// [tests/cases/compiler/declFilePrivateMethodOverloads.ts] ////

//// [declFilePrivateMethodOverloads.ts]
interface IContext {
    someMethod();
}
class c1 {
    private _forEachBindingContext(bindingContext: IContext, fn: (bindingContext: IContext) => void);
    private _forEachBindingContext(bindingContextArray: Array<IContext>, fn: (bindingContext: IContext) => void);
    private _forEachBindingContext(context, fn: (bindingContext: IContext) => void): void {
        // Function here
    }

    private overloadWithArityDifference(bindingContext: IContext);
    private overloadWithArityDifference(bindingContextArray: Array<IContext>, fn: (bindingContext: IContext) => void);
    private overloadWithArityDifference(context): void {
        // Function here
    }
}
declare class c2 {
    private overload1(context, fn);

    private overload2(context);
    private overload2(context, fn);
}

//// [declFilePrivateMethodOverloads.js]
"use strict";
class c1 {
    _forEachBindingContext(context, fn) {
        // Function here
    }
    overloadWithArityDifference(context) {
        // Function here
    }
}


//// [declFilePrivateMethodOverloads.d.ts]
interface IContext {
    someMethod(): any;
}
class c1 {
    private _forEachBindingContext;
    private overloadWithArityDifference;
}
class c2 {
    private overload1;
    private overload2;
}


//// [DtsFileErrors]


declFilePrivateMethodOverloads.d.ts(4,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declFilePrivateMethodOverloads.d.ts (1 errors) ====
    interface IContext {
        someMethod(): any;
    }
    class c1 {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        private _forEachBindingContext;
        private overloadWithArityDifference;
    }
    class c2 {
        private overload1;
        private overload2;
    }
    