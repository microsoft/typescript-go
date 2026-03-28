// @strict: true

// Two overloads with multiple diagnostics per candidate (maxDiagCount > 1)
// This should select only the candidate with fewest diagnostics
function baz(x: string, y: string): string;
function baz(x: number, y: number): number;
function baz(x: any, y: any): any {
    return x;
}

var z = baz(true, true);
