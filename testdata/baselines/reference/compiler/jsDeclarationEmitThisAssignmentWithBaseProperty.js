//// [tests/cases/compiler/jsDeclarationEmitThisAssignmentWithBaseProperty.ts] ////

//// [component.d.ts]
export class Component {
    state: any;
    constructor(props?: any);
}

export class WithAccessor {
    get value(): number;
    set value(v: number);
}

//// [main.js]
import { Component, WithAccessor } from "./component";

export class C1 extends Component {
    state = { count: 0 };
}

export class C2 extends Component {
    constructor() {
        super({});
        this.state = { count: 0 };
    }
}

export class C3 extends Component {
    update() {
        this.state = { count: 1 };
    }
}

export class C4 extends WithAccessor {
    constructor() {
        super();
        this.value = 1;
    }
}

/** @implements {WithAccessor} */
export class C5 {
    constructor() {
        this.value = 1;
    }
}

//// [mainTs.ts]
import { Component } from "./component";

export class C1 extends Component {
    state = { count: 0 };
}

export class C2 extends Component {
    constructor() {
        super({});
        this.state = { count: 0 };
    }
}




//// [main.d.ts]
import { Component, WithAccessor } from "./component";
export declare class C1 extends Component {
    state: {
        count: number;
    };
}
export declare class C2 extends Component {
    state: {
        count: number;
    };
    constructor();
}
export declare class C3 extends Component {
    update(): void;
}
export declare class C4 extends WithAccessor {
    constructor();
}
/** @implements {WithAccessor} */
export declare class C5 implements WithAccessor {
    constructor();
}
//// [mainTs.d.ts]
import { Component } from "./component";
export declare class C1 extends Component {
    state: {
        count: number;
    };
}
export declare class C2 extends Component {
    constructor();
}


//// [DtsFileErrors]


main.d.ts(20,22): error TS2720: Class 'C5' incorrectly implements class 'WithAccessor'. Did you mean to extend 'WithAccessor' and inherit its members as a subclass?
  Property 'value' is missing in type 'C5' but required in type 'WithAccessor'.


==== component.d.ts (0 errors) ====
    export class Component {
        state: any;
        constructor(props?: any);
    }
    
    export class WithAccessor {
        get value(): number;
        set value(v: number);
    }
    
==== main.d.ts (1 errors) ====
    import { Component, WithAccessor } from "./component";
    export declare class C1 extends Component {
        state: {
            count: number;
        };
    }
    export declare class C2 extends Component {
        state: {
            count: number;
        };
        constructor();
    }
    export declare class C3 extends Component {
        update(): void;
    }
    export declare class C4 extends WithAccessor {
        constructor();
    }
    /** @implements {WithAccessor} */
    export declare class C5 implements WithAccessor {
                         ~~
!!! error TS2720: Class 'C5' incorrectly implements class 'WithAccessor'. Did you mean to extend 'WithAccessor' and inherit its members as a subclass?
!!! error TS2720:   Property 'value' is missing in type 'C5' but required in type 'WithAccessor'.
!!! related TS2728 component.d.ts:7:9: 'value' is declared here.
        constructor();
    }
    
==== mainTs.d.ts (0 errors) ====
    import { Component } from "./component";
    export declare class C1 extends Component {
        state: {
            count: number;
        };
    }
    export declare class C2 extends Component {
        constructor();
    }
    