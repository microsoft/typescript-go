// @target: esnext
// @experimentalDecorators: true
// @strict: true

// Repro from issue: class name used as object literal key should not be renamed to the alias

const dec = (t: any) => t;

@dec
class SessionAuth {
    static requirement() {
        return { SessionAuth: [] };
    }
}

// Method shorthand name in object literal should not be renamed
@dec
class Foo {
    static methods() {
        return {
            Foo() { return 1; }
        };
    }
}

// Class field name should not be renamed
@dec
class Bar {
    Bar = 42;
}

// Object binding pattern property key should not be renamed
@dec
class Baz {
    static unwrap(value: any) {
        const { Baz: renamed } = value;
        return renamed;
    }
}

// Self-reference in expression positions SHOULD still be renamed
@dec
class SelfRef {
    static instance = new SelfRef();
    method() {
        return SelfRef;
    }
}
