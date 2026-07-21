// @noImplicitReturns: true
// @strict: true

// Exhaustiveness helper that returns never
export function assertUnreachable(x: never): never {
    throw new Error('unreachable: ' + x);
}

// Test 1: Simple case - call never-returning function in default
export function pick(key: 'a' | 'b' | 'c') {
    switch (key) {
        case 'a':
            return 1;
        case 'b':
            return 2;
        case 'c':
            return 3;
        default:
            assertUnreachable(key);
    }
}

// Test 2: With explicit return type annotation
export function pick2(key: 'a' | 'b' | 'c'): number {
    switch (key) {
        case 'a':
            return 1;
        case 'b':
            return 2;
        case 'c':
            return 3;
        default:
            assertUnreachable(key);
    }
}

// Test 3: Using a method that returns never (dotted name call)
class Assertions {
    fail(x: never): never {
        throw new Error('unreachable');
    }
}

export function pick3(key: 'a' | 'b' | 'c', assertions: Assertions) {
    switch (key) {
        case 'a':
            return 1;
        case 'b':
            return 2;
        case 'c':
            return 3;
        default:
            assertions.fail(key);
    }
}
