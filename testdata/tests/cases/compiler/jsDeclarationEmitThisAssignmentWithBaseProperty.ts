// @allowJs: true
// @checkJs: true
// @declaration: true
// @emitDeclarationOnly: true
// @target: esnext

// @filename: component.d.ts
export class Component {
    state: any;
    constructor(props?: any);
}

export class WithAccessor {
    get value(): number;
    set value(v: number);
}

// @filename: main.js
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

// @filename: mainTs.ts
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
