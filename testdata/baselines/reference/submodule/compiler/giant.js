//// [tests/cases/compiler/giant.ts] ////

//// [giant.ts]
/*
    Prefixes
    p -> public
    r -> private
    i -> import
    e -> export
    a -> ambient
    t -> static
    s -> set
    g -> get

    MAX DEPTH 3 LEVELS
*/
const p = "propName";
var V;
function F() { };
class C {
    constructor () { }
    public pV;
    private rV;
    public pF() { }
    private rF() { }
    public pgF() { }
    public get pgF()
    public psF(param:any) { }
    public set psF(param:any)
    private rgF() { }
    private get rgF()
    private rsF(param:any) { }
    private set rsF(param:any)
    static tV;
    static tF() { }
    static tsF(param:any) { }
    static set tsF(param:any)
    static tgF() { }
    static get tgF()
}
interface I {
    //Call Signature
    ();
    (): number;
    (p);
    (p1: string);
    (p2?: string);
    (...p3: any[]);
    (p4: string, p5?: string);
    (p6: string, ...p7: any[]);
    //(p8?: string, ...p9: any[]);
    //(p10:string, p8?: string, ...p9: any[]);
    
    //Construct Signature
    new ();
    new (): number;
    new (p: string);
    new (p2?: string);
    new (...p3: any[]);
    new (p4: string, p5?: string);
    new (p6: string, ...p7: any[]);

    //Index Signature
    [p];
    [p1: string];
    [p2: string, p3: number];

    //Property Signature
    p;
    p1?;
    p2?: string;
    
    //Function Signature
    p3();
    p4? ();
    p5? (): void;
    p6(pa1): void;
    p7(pa1, pa2): void;
    p7? (pa1, pa2): void;
}
module M {
    var V;
    function F() { };
    class C {
        constructor () { }
        public pV;
        private rV;
        public pF() { }
        private rF() { }
        public pgF() { }
        public get pgF()
        public psF(param:any) { }
        public set psF(param:any)
        private rgF() { }
        private get rgF()
        private rsF(param:any) { }
        private set rsF(param:any)
        static tV;
        static tF() { }
        static tsF(param:any) { }
        static set tsF(param:any)
        static tgF() { }
        static get tgF()
    }
    interface I {
        //Call Signature
        ();
        (): number;
        (p);
        (p1: string);
        (p2?: string);
        (...p3: any[]);
        (p4: string, p5?: string);
        (p6: string, ...p7: any[]);
        //(p8?: string, ...p9: any[]);
        //(p10:string, p8?: string, ...p9: any[]);
    
        //Construct Signature
        new ();
        new (): number;
        new (p: string);
        new (p2?: string);
        new (...p3: any[]);
        new (p4: string, p5?: string);
        new (p6: string, ...p7: any[]);

        //Index Signature
        [p];
        [p1: string];
        [p2: string, p3: number];

        //Property Signature
        p;
        p1?;
        p2?: string;
    
        //Function Signature
        p3();
        p4? ();
        p5? (): void;
        p6(pa1): void;
        p7(pa1, pa2): void;
        p7? (pa1, pa2): void;
    }
    module M {
        var V;
        function F() { };
        class C { };
        interface I { };
        module M { };
        export var eV;
        export function eF() { };
        export class eC { };
        export interface eI { };
        export module eM { };
        export declare var eaV;
        export declare function eaF() { };
        export declare class eaC { };
        export declare module eaM { };
    }
    export var eV;
    export function eF() { };
    export class eC {
        constructor () { }
        public pV;
        private rV;
        public pF() { }
        private rF() { }
        public pgF() { }
        public get pgF()
        public psF(param:any) { }
        public set psF(param:any)
        private rgF() { }
        private get rgF()
        private rsF(param:any) { }
        private set rsF(param:any)
        static tV;
        static tF() { }
        static tsF(param:any) { }
        static set tsF(param:any)
        static tgF() { }
        static get tgF()
    }
    export interface eI {
        //Call Signature
        ();
        (): number;
        (p);
        (p1: string);
        (p2?: string);
        (...p3: any[]);
        (p4: string, p5?: string);
        (p6: string, ...p7: any[]);
        //(p8?: string, ...p9: any[]);
        //(p10:string, p8?: string, ...p9: any[]);
    
        //Construct Signature
        new ();
        new (): number;
        new (p: string);
        new (p2?: string);
        new (...p3: any[]);
        new (p4: string, p5?: string);
        new (p6: string, ...p7: any[]);

        //Index Signature
        [p];
        [p1: string];
        [p2: string, p3: number];

        //Property Signature
        p;
        p1?;
        p2?: string;
    
        //Function Signature
        p3();
        p4? ();
        p5? (): void;
        p6(pa1): void;
        p7(pa1, pa2): void;
        p7? (pa1, pa2): void;
    }
    export module eM {
        var V;
        function F() { };
        class C { };
        interface I { };
        module M { };
        export var eV;
        export function eF() { };
        export class eC { };
        export interface eI { };
        export module eM { };
        export declare var eaV;
        export declare function eaF() { };
        export declare class eaC { };
        export declare module eaM { };
    }
    export declare var eaV;
    export declare function eaF() { };
    export declare class eaC {
        constructor () { }
        public pV;
        private rV;
        public pF() { }
        private rF() { }
        public pgF() { }
        public get pgF()
        public psF(param:any) { }
        public set psF(param:any)
        private rgF() { }
        private get rgF()
        private rsF(param:any) { }
        private set rsF(param:any)
        static tV;
        static tF() { }
        static tsF(param:any) { }
        static set tsF(param:any)
        static tgF() { }
        static get tgF()
    }
    export declare module eaM {
        var V;
        function F() { };
        class C { }
        interface I { }
        module M { }
        export var eV;
        export function eF() { };
        export class eC { }
        export interface eI { }
        export module eM { }
    }
}
export var eV;
export function eF() { };
export class eC {
    constructor () { }
    public pV;
    private rV;
    public pF() { }
    private rF() { }
    public pgF() { }
    public get pgF()
    public psF(param:any) { }
    public set psF(param:any)
    private rgF() { }
    private get rgF()
    private rsF(param:any) { }
    private set rsF(param:any)
    static tV;
    static tF() { }
    static tsF(param:any) { }
    static set tsF(param:any)
    static tgF() { }
    static get tgF()
}
export interface eI {
    //Call Signature
    ();
    (): number;
    (p);
    (p1: string);
    (p2?: string);
    (...p3: any[]);
    (p4: string, p5?: string);
    (p6: string, ...p7: any[]);
    //(p8?: string, ...p9: any[]);
    //(p10:string, p8?: string, ...p9: any[]);
    
    //Construct Signature
    new ();
    new (): number;
    new (p: string);
    new (p2?: string);
    new (...p3: any[]);
    new (p4: string, p5?: string);
    new (p6: string, ...p7: any[]);

    //Index Signature
    [p];
    [p1: string];
    [p2: string, p3: number];

    //Property Signature
    p;
    p1?;
    p2?: string;
    
    //Function Signature
    p3();
    p4? ();
    p5? (): void;
    p6(pa1): void;
    p7(pa1, pa2): void;
    p7? (pa1, pa2): void;
}
export module eM {
    var V;
    function F() { };
    class C {
        constructor () { }
        public pV;
        private rV;
        public pF() { }
        private rF() { }
        public pgF() { }
        public get pgF()
        public psF(param:any) { }
        public set psF(param:any)
        private rgF() { }
        private get rgF()
        private rsF(param:any) { }
        private set rsF(param:any)
        static tV;
        static tF() { }
        static tsF(param:any) { }
        static set tsF(param:any)
        static tgF() { }
        static get tgF()
    }
    interface I {
        //Call Signature
        ();
        (): number;
        (p);
        (p1: string);
        (p2?: string);
        (...p3: any[]);
        (p4: string, p5?: string);
        (p6: string, ...p7: any[]);
        //(p8?: string, ...p9: any[]);
        //(p10:string, p8?: string, ...p9: any[]);
    
        //Construct Signature
        new ();
        new (): number;
        new (p: string);
        new (p2?: string);
        new (...p3: any[]);
        new (p4: string, p5?: string);
        new (p6: string, ...p7: any[]);

        //Index Signature
        [p];
        [p1: string];
        [p2: string, p3: number];

        //Property Signature
        p;
        p1?;
        p2?: string;
    
        //Function Signature
        p3();
        p4? ();
        p5? (): void;
        p6(pa1): void;
        p7(pa1, pa2): void;
        p7? (pa1, pa2): void;
    }
    module M {
        var V;
        function F() { };
        class C { };
        interface I { };
        module M { };
        export var eV;
        export function eF() { };
        export class eC { };
        export interface eI { };
        export module eM { };
        export declare var eaV;
        export declare function eaF() { };
        export declare class eaC { };
        export declare module eaM { };
    }
    export var eV;
    export function eF() { };
    export class eC {
        constructor () { }
        public pV;
        private rV;
        public pF() { }
        private rF() { }
        public pgF() { }
        public get pgF()
        public psF(param:any) { }
        public set psF(param:any)
        private rgF() { }
        private get rgF()
        private rsF(param:any) { }
        private set rsF(param:any)
        static tV;
        static tF() { }
        static tsF(param:any) { }
        static set tsF(param:any)
        static tgF() { }
        static get tgF()
    }
    export interface eI {
        //Call Signature
        ();
        (): number;
        (p);
        (p1: string);
        (p2?: string);
        (...p3: any[]);
        (p4: string, p5?: string);
        (p6: string, ...p7: any[]);
        //(p8?: string, ...p9: any[]);
        //(p10:string, p8?: string, ...p9: any[]);
    
        //Construct Signature
        new ();
        new (): number;
        new (p: string);
        new (p2?: string);
        new (...p3: any[]);
        new (p4: string, p5?: string);
        new (p6: string, ...p7: any[]);

        //Index Signature
        [p];
        [p1: string];
        [p2: string, p3: number];

        //Property Signature
        p;
        p1?;
        p2?: string;
    
        //Function Signature
        p3();
        p4? ();
        p5? (): void;
        p6(pa1): void;
        p7(pa1, pa2): void;
        p7? (pa1, pa2): void;
    }
    export module eM {
        var V;
        function F() { };
        class C { };
        interface I { };
        module M { };
        export var eV;
        export function eF() { };
        export class eC { };
        export interface eI { };
        export module eM { };
        export declare var eaV;
        export declare function eaF() { };
        export declare class eaC { };
        export declare module eaM { };
    }
    export declare var eaV;
    export declare function eaF() { };
    export declare class eaC {
        constructor () { }
        public pV;
        private rV;
        public pF() { }
        private rF() { }
        public pgF() { }
        public get pgF()
        public psF(param:any) { }
        public set psF(param:any)
        private rgF() { }
        private get rgF()
        private rsF(param:any) { }
        private set rsF(param:any)
        static tV;
        static tF() { }
        static tsF(param:any) { }
        static set tsF(param:any)
        static tgF() { }
        static get tgF()
    }
    export declare module eaM {
        var V;
        function F() { };
        class C { }
        interface I { }
        module M { }
        export var eV;
        export function eF() { };
        export class eC { }
        export interface eI { }
        export module eM { }
    }
}
export declare var eaV;
export declare function eaF() { };
export declare class eaC {
    constructor () { }
    public pV;
    private rV;
    public pF() { }
    private rF() { }
    public pgF() { }
    public get pgF()
    public psF(param:any) { }
    public set psF(param:any)
    private rgF() { }
    private get rgF()
    private rsF(param:any) { }
    private set rsF(param:any)
    static tV;
    static tF() { }
    static tsF(param:any) { }
    static set tsF(param:any)
    static tgF() { }
    static get tgF()
}
export declare module eaM {
    var V;
    function F() { };
    class C {
        constructor () { }
        public pV;
        private rV;
        public pF() { }
        static tV;
        static tF() { }
    }
    interface I {
        //Call Signature
        ();
        (): number;
        (p: string);
        (p2?: string);
        (...p3: any[]);
        (p4: string, p5?: string);
        (p6: string, ...p7: any[]);
        //(p8?: string, ...p9: any[]);
        //(p10:string, p8?: string, ...p9: any[]);
    
        //Construct Signature
        new ();
        new (): number;
        new (p: string);
        new (p2?: string);
        new (...p3: any[]);
        new (p4: string, p5?: string);
        new (p6: string, ...p7: any[]);

        //Index Signature
        [p];
        [p1: string];
        [p2: string, p3: number];

        //Property Signature
        p;
        p1?;
        p2?: string;
    
        //Function Signature
        p3();
        p4? ();
        p5? (): void;
        p6(pa1): void;
        p7(pa1, pa2): void;
        p7? (pa1, pa2): void;
    }
    module M {
        var V;
        function F() { };
        class C { }
        interface I { }
        module M { }
        export var eV;
        export function eF() { };
        export class eC { }
        export interface eI { }
        export module eM { }
        export declare var eaV
        export declare function eaF() { };
        export declare class eaC { }
        export declare module eaM { }
    }
    export var eV;
    export function eF() { };
    export class eC {
        constructor () { }
        public pV;
        private rV;
        public pF() { }
        static tV
        static tF() { }
    }
    export interface eI {
        //Call Signature
        ();
        (): number;
        (p);
        (p1: string);
        (p2?: string);
        (...p3: any[]);
        (p4: string, p5?: string);
        (p6: string, ...p7: any[]);
        //(p8?: string, ...p9: any[]);
        //(p10:string, p8?: string, ...p9: any[]);
    
        //Construct Signature
        new ();
        new (): number;
        new (p: string);
        new (p2?: string);
        new (...p3: any[]);
        new (p4: string, p5?: string);
        new (p6: string, ...p7: any[]);

        //Index Signature
        [p];
        [p1: string];
        [p2: string, p3: number];

        //Property Signature
        p;
        p1?;
        p2?: string;
    
        //Function Signature
        p3();
        p4? ();
        p5? (): void;
        p6(pa1): void;
        p7(pa1, pa2): void;
        p7? (pa1, pa2): void;
    }
    export module eM {
        var V;
        function F() { };
        class C { }
        module M { }
        export var eV;
        export function eF() { };
        export class eC { }
        export interface eI { }
        export module eM { }
    }
}

//// [giant.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.eM = exports.eC = exports.eV = void 0;
exports.eF = eF;
/*
    Prefixes
    p -> public
    r -> private
    i -> import
    e -> export
    a -> ambient
    t -> static
    s -> set
    g -> get

    MAX DEPTH 3 LEVELS
*/
const p = "propName";
var V;
function F() { }
;
class C {
    constructor() { }
    pV;
    rV;
    pF() { }
    rF() { }
    pgF() { }
    psF(param) { }
    rgF() { }
    rsF(param) { }
    static tV;
    static tF() { }
    static tsF(param) { }
    static tgF() { }
}
var M;
(function (M_1) {
    var V;
    function F() { }
    ;
    class C {
        constructor() { }
        pV;
        rV;
        pF() { }
        rF() { }
        pgF() { }
        psF(param) { }
        rgF() { }
        rsF(param) { }
        static tV;
        static tF() { }
        static tsF(param) { }
        static tgF() { }
    }
    let M;
    (function (M) {
        var V;
        function F() { }
        ;
        class C {
        }
        ;
        ;
        ;
        function eF() { }
        M.eF = eF;
        ;
        class eC {
        }
        M.eC = eC;
        ;
        ;
        ;
        ;
        ;
        ;
    })(M || (M = {}));
    function eF() { }
    M_1.eF = eF;
    ;
    class eC {
        constructor() { }
        pV;
        rV;
        pF() { }
        rF() { }
        pgF() { }
        psF(param) { }
        rgF() { }
        rsF(param) { }
        static tV;
        static tF() { }
        static tsF(param) { }
        static tgF() { }
    }
    M_1.eC = eC;
    let eM;
    (function (eM) {
        var V;
        function F() { }
        ;
        class C {
        }
        ;
        ;
        ;
        function eF() { }
        eM.eF = eF;
        ;
        class eC {
        }
        eM.eC = eC;
        ;
        ;
        ;
        ;
        ;
        ;
    })(eM = M_1.eM || (M_1.eM = {}));
    ;
})(M || (M = {}));
function eF() { }
;
class eC {
    constructor() { }
    pV;
    rV;
    pF() { }
    rF() { }
    pgF() { }
    psF(param) { }
    rgF() { }
    rsF(param) { }
    static tV;
    static tF() { }
    static tsF(param) { }
    static tgF() { }
}
exports.eC = eC;
var eM;
(function (eM_1) {
    var V;
    function F() { }
    ;
    class C {
        constructor() { }
        pV;
        rV;
        pF() { }
        rF() { }
        pgF() { }
        psF(param) { }
        rgF() { }
        rsF(param) { }
        static tV;
        static tF() { }
        static tsF(param) { }
        static tgF() { }
    }
    let M;
    (function (M) {
        var V;
        function F() { }
        ;
        class C {
        }
        ;
        ;
        ;
        function eF() { }
        M.eF = eF;
        ;
        class eC {
        }
        M.eC = eC;
        ;
        ;
        ;
        ;
        ;
        ;
    })(M || (M = {}));
    function eF() { }
    eM_1.eF = eF;
    ;
    class eC {
        constructor() { }
        pV;
        rV;
        pF() { }
        rF() { }
        pgF() { }
        psF(param) { }
        rgF() { }
        rsF(param) { }
        static tV;
        static tF() { }
        static tsF(param) { }
        static tgF() { }
    }
    eM_1.eC = eC;
    let eM;
    (function (eM) {
        var V;
        function F() { }
        ;
        class C {
        }
        ;
        ;
        ;
        function eF() { }
        eM.eF = eF;
        ;
        class eC {
        }
        eM.eC = eC;
        ;
        ;
        ;
        ;
        ;
        ;
    })(eM = eM_1.eM || (eM_1.eM = {}));
    ;
})(eM || (exports.eM = eM = {}));
;


//// [giant.d.ts]
/*
    Prefixes
    p -> public
    r -> private
    i -> import
    e -> export
    a -> ambient
    t -> static
    s -> set
    g -> get

    MAX DEPTH 3 LEVELS
*/
declare const p = "propName";
export declare var eV: any;
export declare function eF(): void;
export declare class eC {
    constructor();
    pV: any;
    private rV;
    pF(): void;
    private rF;
    pgF(): void;
    get pgF(): any;
    psF(param: any): void;
    set psF(param: any);
    private rgF;
    private get rgF();
    private rsF;
    private set rsF(value);
    static tV: any;
    static tF(): void;
    static tsF(param: any): void;
    static set tsF(param: any);
    static tgF(): void;
    static get tgF(): any;
}
export interface eI {
    //Call Signature
    (): any;
    (): number;
    (p: any): any;
    (p1: string): any;
    (p2?: string): any;
    (...p3: any[]): any;
    (p4: string, p5?: string): any;
    (p6: string, ...p7: any[]): any;
    //(p8?: string, ...p9: any[]);
    //(p10:string, p8?: string, ...p9: any[]);
    //Construct Signature
    new ();
    new ();
    new (p: string);
    new (p2?: string);
    new (...p3: any[]);
    new (p4: string, p5?: string);
    new (p6: string, ...p7: any[]);
    //Index Signature
    [p]: any;
    [p1: string]: any;
    [p2: string, p3: number]: any;
    //Property Signature
    p: any;
    p1?: any;
    p2?: string;
    //Function Signature
    p3(): any;
    p4?(): any;
    p5?(): void;
    p6(pa1: any): void;
    p7(pa1: any, pa2: any): void;
    p7?(pa1: any, pa2: any): void;
}
export declare namespace eM {
    var eV: any;
    function eF(): void;
    class eC {
        constructor();
        pV: any;
        private rV;
        pF(): void;
        private rF;
        pgF(): void;
        get pgF(): any;
        psF(param: any): void;
        set psF(param: any);
        private rgF;
        private get rgF();
        private rsF;
        private set rsF(value);
        static tV: any;
        static tF(): void;
        static tsF(param: any): void;
        static set tsF(param: any);
        static tgF(): void;
        static get tgF(): any;
    }
    interface eI {
        //Call Signature
        (): any;
        (): number;
        (p: any): any;
        (p1: string): any;
        (p2?: string): any;
        (...p3: any[]): any;
        (p4: string, p5?: string): any;
        (p6: string, ...p7: any[]): any;
        //(p8?: string, ...p9: any[]);
        //(p10:string, p8?: string, ...p9: any[]);
        //Construct Signature
        new ();
        new ();
        new (p: string);
        new (p2?: string);
        new (...p3: any[]);
        new (p4: string, p5?: string);
        new (p6: string, ...p7: any[]);
        //Index Signature
        [p]: any;
        [p1: string]: any;
        [p2: string, p3: number]: any;
        //Property Signature
        p: any;
        p1?: any;
        p2?: string;
        //Function Signature
        p3(): any;
        p4?(): any;
        p5?(): void;
        p6(pa1: any): void;
        p7(pa1: any, pa2: any): void;
        p7?(pa1: any, pa2: any): void;
    }
    namespace eM {
        var eV: any;
        function eF(): void;
        class eC {
        }
        interface eI {
        }
        namespace eM { }
        var eaV: any;
        function eaF(): void;
        class eaC {
        }
        namespace eaM { }
    }
    var eaV: any;
    function eaF(): void;
    class eaC {
        constructor();
        pV: any;
        private rV;
        pF(): void;
        private rF;
        pgF(): void;
        get pgF(): any;
        psF(param: any): void;
        set psF(param: any);
        private rgF;
        private get rgF();
        private rsF;
        private set rsF(value);
        static tV: any;
        static tF(): void;
        static tsF(param: any): void;
        static set tsF(param: any);
        static tgF(): void;
        static get tgF(): any;
    }
    namespace eaM {
        var V: any;
        function F(): void;
        class C {
        }
        interface I {
        }
        namespace M { }
        var eV: any;
        function eF(): void;
        class eC {
        }
        interface eI {
        }
        namespace eM { }
    }
}
export declare var eaV: any;
export declare function eaF(): void;
export declare class eaC {
    constructor();
    pV: any;
    private rV;
    pF(): void;
    private rF;
    pgF(): void;
    get pgF(): any;
    psF(param: any): void;
    set psF(param: any);
    private rgF;
    private get rgF();
    private rsF;
    private set rsF(value);
    static tV: any;
    static tF(): void;
    static tsF(param: any): void;
    static set tsF(param: any);
    static tgF(): void;
    static get tgF(): any;
}
export declare namespace eaM {
    var V: any;
    function F(): void;
    class C {
        constructor();
        pV: any;
        private rV;
        pF(): void;
        static tV: any;
        static tF(): void;
    }
    interface I {
        //Call Signature
        (): any;
        (): number;
        (p: string): any;
        (p2?: string): any;
        (...p3: any[]): any;
        (p4: string, p5?: string): any;
        (p6: string, ...p7: any[]): any;
        //(p8?: string, ...p9: any[]);
        //(p10:string, p8?: string, ...p9: any[]);
        //Construct Signature
        new ();
        new ();
        new (p: string);
        new (p2?: string);
        new (...p3: any[]);
        new (p4: string, p5?: string);
        new (p6: string, ...p7: any[]);
        //Index Signature
        [p]: any;
        [p1: string]: any;
        [p2: string, p3: number]: any;
        //Property Signature
        p: any;
        p1?: any;
        p2?: string;
        //Function Signature
        p3(): any;
        p4?(): any;
        p5?(): void;
        p6(pa1: any): void;
        p7(pa1: any, pa2: any): void;
        p7?(pa1: any, pa2: any): void;
    }
    namespace M {
        var V: any;
        function F(): void;
        class C {
        }
        interface I {
        }
        namespace M { }
        var eV: any;
        function eF(): void;
        class eC {
        }
        interface eI {
        }
        namespace eM { }
        var eaV: any;
        function eaF(): void;
        class eaC {
        }
        namespace eaM { }
    }
    var eV: any;
    function eF(): void;
    class eC {
        constructor();
        pV: any;
        private rV;
        pF(): void;
        static tV: any;
        static tF(): void;
    }
    interface eI {
        //Call Signature
        (): any;
        (): number;
        (p: any): any;
        (p1: string): any;
        (p2?: string): any;
        (...p3: any[]): any;
        (p4: string, p5?: string): any;
        (p6: string, ...p7: any[]): any;
        //(p8?: string, ...p9: any[]);
        //(p10:string, p8?: string, ...p9: any[]);
        //Construct Signature
        new ();
        new ();
        new (p: string);
        new (p2?: string);
        new (...p3: any[]);
        new (p4: string, p5?: string);
        new (p6: string, ...p7: any[]);
        //Index Signature
        [p]: any;
        [p1: string]: any;
        [p2: string, p3: number]: any;
        //Property Signature
        p: any;
        p1?: any;
        p2?: string;
        //Function Signature
        p3(): any;
        p4?(): any;
        p5?(): void;
        p6(pa1: any): void;
        p7(pa1: any, pa2: any): void;
        p7?(pa1: any, pa2: any): void;
    }
    namespace eM {
        var V: any;
        function F(): void;
        class C {
        }
        namespace M { }
        var eV: any;
        function eF(): void;
        class eC {
        }
        interface eI {
        }
        namespace eM { }
    }
}
export {};
