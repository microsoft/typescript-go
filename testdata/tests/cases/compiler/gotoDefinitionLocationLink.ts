// @target: esnext
// @strict: true

// @filename: file1.ts
export interface Person {
    name: string;
    age: number;
}

export function createPerson(name: string, age: number): Person {
    return { name, age };
}

// @filename: file2.ts
import { Person, createPerson } from "./file1";

const person: Person = createPerson("John", 30);
console.log(person.name);