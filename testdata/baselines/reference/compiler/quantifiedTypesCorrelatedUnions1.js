//// [tests/cases/compiler/quantifiedTypesCorrelatedUnions1.ts] ////

//// [quantifiedTypesCorrelatedUnions1.ts]
// https://github.com/microsoft/TypeScript/issues/30581#issuecomment-493492463

declare const str: string;
declare const num: number;
function acceptString(str: string) { }
function acceptNumber(num: number) { }

const arr: (<T> [T, (t: NoInfer<T>) => void])[] = [
    [str, acceptString],
    [num, acceptNumber],
    [str, acceptNumber], // error as expected
];

for (const pair of arr) {
    const [arg, func] = pair; // no error
    func(arg);
}


//// [quantifiedTypesCorrelatedUnions1.js]
// https://github.com/microsoft/TypeScript/issues/30581#issuecomment-493492463
function acceptString(str) { }
function acceptNumber(num) { }
const arr = [
    [str, acceptString],
    [num, acceptNumber],
    [str, acceptNumber], // error as expected
];
for (const pair of arr) {
    const [arg, func] = pair; // no error
    func(arg);
}
