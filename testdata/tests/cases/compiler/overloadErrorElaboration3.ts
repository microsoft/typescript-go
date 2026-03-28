// @strict: true

// Three overloads — should show per-overload errors
function bar(x: string): string;
function bar(x: number): number;
function bar(x: boolean): boolean;
function bar(x: any): any {
    return x;
}

var y = bar({ a: 1 });
