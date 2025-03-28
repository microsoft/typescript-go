//// [tests/cases/conformance/classes/members/privateNames/privateNameFieldUnaryMutation.ts] ////

//// [privateNameFieldUnaryMutation.ts]
class C {
    #test: number = 24;
    constructor() {
        this.#test++;
        this.#test--;
        ++this.#test;
        --this.#test;
        const a = this.#test++;
        const b = this.#test--;
        const c = ++this.#test;
        const d = --this.#test;
        for (this.#test = 0; this.#test < 10; ++this.#test) {}
        for (this.#test = 0; this.#test < 10; this.#test++) {}

        (this.#test)++;
        (this.#test)--;
        ++(this.#test);
        --(this.#test);
        const e = (this.#test)++;
        const f = (this.#test)--;
        const g = ++(this.#test);
        const h = --(this.#test);
        for (this.#test = 0; this.#test < 10; ++(this.#test)) {}
        for (this.#test = 0; this.#test < 10; (this.#test)++) {}
    }
    test() {
        this.getInstance().#test++;
        this.getInstance().#test--;
        ++this.getInstance().#test;
        --this.getInstance().#test;
        const a = this.getInstance().#test++;
        const b = this.getInstance().#test--;
        const c = ++this.getInstance().#test;
        const d = --this.getInstance().#test;
        for (this.getInstance().#test = 0; this.getInstance().#test < 10; ++this.getInstance().#test) {}
        for (this.getInstance().#test = 0; this.getInstance().#test < 10; this.getInstance().#test++) {}

        (this.getInstance().#test)++;
        (this.getInstance().#test)--;
        ++(this.getInstance().#test);
        --(this.getInstance().#test);
        const e = (this.getInstance().#test)++;
        const f = (this.getInstance().#test)--;
        const g = ++(this.getInstance().#test);
        const h = --(this.getInstance().#test);
        for (this.getInstance().#test = 0; this.getInstance().#test < 10; ++(this.getInstance().#test)) {}
        for (this.getInstance().#test = 0; this.getInstance().#test < 10; (this.getInstance().#test)++) {}
    }
    getInstance() { return new C(); }
}


//// [privateNameFieldUnaryMutation.js]
class C {
    #test = 24;
    constructor() {
        this.#test++;
        this.#test--;
        ++this.#test;
        --this.#test;
        const a = this.#test++;
        const b = this.#test--;
        const c = ++this.#test;
        const d = --this.#test;
        for (this.#test = 0; this.#test < 10; ++this.#test) { }
        for (this.#test = 0; this.#test < 10; this.#test++) { }
        (this.#test)++;
        (this.#test)--;
        ++(this.#test);
        --(this.#test);
        const e = (this.#test)++;
        const f = (this.#test)--;
        const g = ++(this.#test);
        const h = --(this.#test);
        for (this.#test = 0; this.#test < 10; ++(this.#test)) { }
        for (this.#test = 0; this.#test < 10; (this.#test)++) { }
    }
    test() {
        this.getInstance().#test++;
        this.getInstance().#test--;
        ++this.getInstance().#test;
        --this.getInstance().#test;
        const a = this.getInstance().#test++;
        const b = this.getInstance().#test--;
        const c = ++this.getInstance().#test;
        const d = --this.getInstance().#test;
        for (this.getInstance().#test = 0; this.getInstance().#test < 10; ++this.getInstance().#test) { }
        for (this.getInstance().#test = 0; this.getInstance().#test < 10; this.getInstance().#test++) { }
        (this.getInstance().#test)++;
        (this.getInstance().#test)--;
        ++(this.getInstance().#test);
        --(this.getInstance().#test);
        const e = (this.getInstance().#test)++;
        const f = (this.getInstance().#test)--;
        const g = ++(this.getInstance().#test);
        const h = --(this.getInstance().#test);
        for (this.getInstance().#test = 0; this.getInstance().#test < 10; ++(this.getInstance().#test)) { }
        for (this.getInstance().#test = 0; this.getInstance().#test < 10; (this.getInstance().#test)++) { }
    }
    getInstance() { return new C(); }
}
