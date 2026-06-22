//// [tests/cases/compiler/jsDeclarationEmitThisAssignmentWithBaseProperty.ts] ////

//// [component.d.ts]
export class Component {
    state: any;
    constructor(props?: any);
}

//// [main.js]
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
import { Component } from "./component";
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
