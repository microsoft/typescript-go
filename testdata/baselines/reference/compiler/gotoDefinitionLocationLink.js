//// [tests/cases/compiler/gotoDefinitionLocationLink.ts] ////

//// [file1.ts]
export interface Person {
    name: string;
    age: number;
}

export function createPerson(name: string, age: number): Person {
    return { name, age };
}

//// [file2.ts]
import { Person, createPerson } from "./file1";

const person: Person = createPerson("John", 30);
console.log(person.name);

//// [file1.js]
export function createPerson(name, age) {
    return { name, age };
}
//// [file2.js]
import { createPerson } from "./file1";
const person = createPerson("John", 30);
console.log(person.name);
