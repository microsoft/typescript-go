// @module: esnext
// @filename: /a.d.ts
declare module "dep" {
    type Arr = string[];
    interface DepModule {
        (): void;
        readonly data: Arr;
    }
    const dep: DepModule;
    export = dep;
}

declare module "wrapper" {
    import * as dep from "dep";
    export { dep as ns };
}

// @filename: /b.d.ts
declare module "wrapper" {
    export const ns: number;
}

// @filename: /main.ts
import * as w from "wrapper";
