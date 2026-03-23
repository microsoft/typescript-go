// @module: commonjs
// @declaration: true
// @filename: a.ts
type Constructor<T = {}> = new (...args: any[]) => T;
interface Printable {
    print(): void;
}
function Mixin<TBase extends Constructor>(Base: TBase) {
    return class extends Base implements Printable {
        print() {}
    };
}
class CoreBase {
    id: number = 0;
    static Printable: Printable = { print() {} };
}
const Mixed = Mixin(CoreBase);
export = Mixed;
export { Printable };
// @filename: b.ts
import Mixed = require("./a");
import { Printable } from "./a";
class App extends Mixed {
    doPrint(p: Printable) {
        p.print();
    }
}
