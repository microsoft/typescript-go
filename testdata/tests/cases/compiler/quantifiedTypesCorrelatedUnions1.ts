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
