// @declaration: true
// @noUnusedLocals: true
// @strict: true
// @target: esnext

export class C1 {
    #a = 1;

    static #b = "";

    #c() {
        return 1 as const;
    }

    get #d() {
        return new Date();
    }

    ["#a"] = true;

    a: this[#a] = 1;

    test1(a: C1) {
        const b: C1[#a] = 1;
        const c: typeof a[#a] = 1;
        const d: C2[#a] = 1;
        const e: typeof C1[#b] = "";
        const f: C1["#a"] = true;
        const g: C1[#c] = this.#c;
        const h: C1[#d] = this.#d;

        type A = #a;
        const i: C1[A] = 1;
        const j: C1[#a | "#a"] = Math.random() ? 1 : true;
        const k: C1[#a] = "";

        type B = C1[#b];
        type C = typeof C1[#a];
        const l = undefined as unknown as B;
        const n = undefined as unknown as C;

        return [b, c, d, e, f, g, h, i, j, k, l, n];
    }

    test2<T extends C1>(a: T): T[#a] {
        const b: T[#a] = a.#a;
        return b;
    }

    test3<T extends C1 | { a: string }>(): T[#a] {
        return undefined as any;
    }

    test4(): any[#a] {
        return "";
    }

    test5(): never[#a] {
        throw new Error();
    }
}

export class C2 {
    #a = 1;

    test1(): void {
        class C {
            test1(a: C2): typeof a[#a] {
                return 1;
            }
        }
        new C().test1(this);
    }
}
