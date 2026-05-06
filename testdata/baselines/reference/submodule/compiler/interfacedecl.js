//// [tests/cases/compiler/interfacedecl.ts] ////

//// [interfacedecl.ts]
interface a0 {
    (): string;
    (a, b, c?: string): number;
    
    new (): string;
    new (s: string);

    [n: number]: ()=>string;
    [s: string]: any;

    p1;
    p2: string;
    p3?;
    p4?: number;
    p5: (s: number) =>string;

    f1();
    f2? ();
    f3(a: string): number;
    f4? (s: number): string;
}


interface a1 {
    [n: number]: number;
}

interface a2 {
    [s: string]: number;
}

interface a {
}

interface b extends a {
}

interface c extends a, b {
}

interface d extends a {
}

class c1 implements a {
}
var instance2 = new c1();

//// [interfacedecl.js]
"use strict";
class c1 {
}
var instance2 = new c1();


//// [interfacedecl.d.ts]
interface a0 {
    (): string;
    (a: any, b: any, c?: string): number;
    new (): string;
    new (s: string): any;
    [n: number]: () => string;
    [s: string]: any;
    p1: any;
    p2: string;
    p3?: any;
    p4?: number;
    p5: (s: number) => string;
    f1(): any;
    f2?(): any;
    f3(a: string): number;
    f4?(s: number): string;
}
interface a1 {
    [n: number]: number;
}
interface a2 {
    [s: string]: number;
}
interface a {
}
interface b extends a {
}
interface c extends a, b {
}
interface d extends a {
}
class c1 implements a {
}
var instance2: c1;


//// [DtsFileErrors]


interfacedecl.d.ts(32,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== interfacedecl.d.ts (1 errors) ====
    interface a0 {
        (): string;
        (a: any, b: any, c?: string): number;
        new (): string;
        new (s: string): any;
        [n: number]: () => string;
        [s: string]: any;
        p1: any;
        p2: string;
        p3?: any;
        p4?: number;
        p5: (s: number) => string;
        f1(): any;
        f2?(): any;
        f3(a: string): number;
        f4?(s: number): string;
    }
    interface a1 {
        [n: number]: number;
    }
    interface a2 {
        [s: string]: number;
    }
    interface a {
    }
    interface b extends a {
    }
    interface c extends a, b {
    }
    interface d extends a {
    }
    class c1 implements a {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    }
    var instance2: c1;
    