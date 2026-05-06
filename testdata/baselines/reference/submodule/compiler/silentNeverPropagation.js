//// [tests/cases/compiler/silentNeverPropagation.ts] ////

//// [silentNeverPropagation.ts]
// Repro from #45041

type ModuleWithState<TState> = {
    state: TState;
};

type State = {
    a: number;
};

type MoreState = {
    z: string;
};

declare function createModule<TState, TActions>(state: TState, actions: TActions): ModuleWithState<TState> & TActions;

declare function convert<TState, TActions>(m: ModuleWithState<TState> & TActions): ModuleWithState<TState & MoreState> & TActions;

const breaks = convert(
    createModule({ a: 12 }, { foo() { return true } })
);

breaks.state.a
breaks.state.z
breaks.foo()


//// [silentNeverPropagation.js]
"use strict";
// Repro from #45041
const breaks = convert(createModule({ a: 12 }, { foo() { return true; } }));
breaks.state.a;
breaks.state.z;
breaks.foo();


//// [silentNeverPropagation.d.ts]
type ModuleWithState<TState> = {
    state: TState;
};
type State = {
    a: number;
};
type MoreState = {
    z: string;
};
function createModule<TState, TActions>(state: TState, actions: TActions): ModuleWithState<TState> & TActions;
function convert<TState, TActions>(m: ModuleWithState<TState> & TActions): ModuleWithState<TState & MoreState> & TActions;
const breaks: ModuleWithState<{
    a: number;
} & MoreState> & ModuleWithState<{
    a: number;
}> & {
    foo(): true;
};


//// [DtsFileErrors]


silentNeverPropagation.d.ts(10,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== silentNeverPropagation.d.ts (1 errors) ====
    type ModuleWithState<TState> = {
        state: TState;
    };
    type State = {
        a: number;
    };
    type MoreState = {
        z: string;
    };
    function createModule<TState, TActions>(state: TState, actions: TActions): ModuleWithState<TState> & TActions;
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    function convert<TState, TActions>(m: ModuleWithState<TState> & TActions): ModuleWithState<TState & MoreState> & TActions;
    const breaks: ModuleWithState<{
        a: number;
    } & MoreState> & ModuleWithState<{
        a: number;
    }> & {
        foo(): true;
    };
    