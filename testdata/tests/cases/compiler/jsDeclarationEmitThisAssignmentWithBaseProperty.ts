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

// @filename: main.js
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
