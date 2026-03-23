// @module: commonjs
// @declaration: true
// @filename: a.d.ts
declare class Lib {
    value: number;
}
declare namespace Lib {
    class Component {
        render(): void;
    }
}
export = Lib;
// @filename: b.ts
import { Component } from "./a";
class App extends Component {
    render() { }
}
