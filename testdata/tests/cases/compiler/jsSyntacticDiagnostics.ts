// @allowJs: true
// @checkJs: true
// @noEmit: true

// @filename: test.js

// Type annotations should be flagged as errors
function func(x: number): string {
    return x.toString();
}

// Interface declarations should be flagged as errors
interface Person {
    name: string;
    age: number;
}

// Type alias declarations should be flagged as errors
type StringOrNumber = string | number;

// Enum declarations should be flagged as errors
enum Color {
    Red,
    Green,
    Blue
}

// Module declarations should be flagged as errors
module MyModule {
    export var x = 1;
}

// Namespace declarations should be flagged as errors
namespace MyNamespace {
    export var y = 2;
}

// Non-null assertions should be flagged as errors
let value = getValue()!;

// Type assertions should be flagged as errors
let result = (value as string).toUpperCase();

// Satisfies expressions should be flagged as errors
let config = {} satisfies Config;

// Import type should be flagged as errors
import type { SomeType } from './other';

// Export type should be flagged as errors
export type { SomeType };

// Import equals should be flagged as errors
import lib = require('./lib');

// Export equals should be flagged as errors
export = MyModule;

// TypeScript modifiers should be flagged as errors
class MyClass {
    public name: string;
    private age: number;
    protected id: number;
    readonly value: number;
    
    constructor(public x: number, private y: number) {
        this.name = '';
        this.age = 0;
        this.id = 0;
        this.value = 0;
    }
}

// Optional parameters should be flagged as errors
function optionalParam(x?: number) {
    return x || 0;
}

// Signature declarations should be flagged as errors
function signatureOnly(x: number): string;

// Type parameters should be flagged as errors
function generic<T>(x: T): T {
    return x;
}

// Type arguments should be flagged as errors
let array = Array<string>();

// Implements clause should be flagged as errors
class MyClassWithImplements implements Person {
    name = '';
    age = 0;
}

function getValue(): any {
    return null;
}

interface Config {
    name: string;
}