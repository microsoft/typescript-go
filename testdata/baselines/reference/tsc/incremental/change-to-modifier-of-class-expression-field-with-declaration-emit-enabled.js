
currentDirectory::/home/src/workspaces/project
useCaseSensitiveFileNames::true
Input::--incremental
//// [/home/src/tslibs/TS/Lib/lib.d.ts] *new* 
/// <reference no-default-lib="true"/>
interface Boolean {}
interface Function {}
interface CallableFunction {}
interface NewableFunction {}
interface IArguments {}
interface Number { toExponential: any; }
interface Object {}
interface RegExp {}
interface String { charAt: any; }
interface Array<T> { length: number; [n: number]: T; }
interface ReadonlyArray<T> {}
interface SymbolConstructor {
    (desc?: string | number): symbol;
    for(name: string): symbol;
    readonly toStringTag: symbol;
}
declare var Symbol: SymbolConstructor;
interface Symbol {
    readonly [Symbol.toStringTag]: string;
}
declare const console: { log(msg: any): void; };
					type ReturnType<T extends (...args: any) => any> = T extends (...args: any) => infer R ? R : any;
                    type InstanceType<T extends abstract new (...args: any) => any> = T extends abstract new (...args: any) => infer R ? R : any;
//// [/home/src/workspaces/project/MessageablePerson.ts] *new* 

                        const Messageable = () => {
                            return class MessageableClass {
                                public message = 'hello';
                            }
                        };
                        const wrapper = () => Messageable();
                        type MessageablePerson = InstanceType<ReturnType<typeof wrapper>>;
                        export default MessageablePerson;
//// [/home/src/workspaces/project/main.ts] *new* 

                        import MessageablePerson from './MessageablePerson.js';
                        function logMessage( person: MessageablePerson ) {
                            console.log( person.message );
                        }
//// [/home/src/workspaces/project/tsconfig.json] *new* 
{ "compilerOptions": { "module": "esnext", "declaration": true  } }

ExitStatus:: 0

CompilerOptions::{
    "incremental": true
}
Output::
//// [/home/src/workspaces/project/MessageablePerson.d.ts] *new* 
declare const wrapper: () => {
    new (): {
        message: string;
    };
};
type MessageablePerson = InstanceType<ReturnType<typeof wrapper>>;
export default MessageablePerson;

//// [/home/src/workspaces/project/MessageablePerson.js] *new* 
const Messageable = () => {
    return class MessageableClass {
        message = 'hello';
    };
};
const wrapper = () => Messageable();
export {};

//// [/home/src/workspaces/project/main.d.ts] *new* 
export {};

//// [/home/src/workspaces/project/main.js] *new* 
function logMessage(person) {
    console.log(person.message);
}
export {};

//// [/home/src/workspaces/project/tsconfig.tsbuildinfo] *new* 
{"version":"FakeTSVersion","fileNames":["../../tslibs/TS/Lib/lib.d.ts","./MessageablePerson.ts","./main.ts"],"fileInfos":[{"version":"575a4e15624573144926595079b1ec30f9c7853bab32f43c0b7db2acfdf038e2","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"ff666de4fdc53b5500de60a9b8c073c9327a9e9326417ef4861b8d2473c7457a","signature":"6ec1f7bdc192ba06258caff3fa202fd577f8f354d676f548500eeb232155cbbe","impliedNodeFormat":1},{"version":"36f0b00de3c707929bf1919e32e5b6053c8730bb00aa779bcdd1925414d68b8c","signature":"8e609bb71c20b858c77f0e9f90bb1319db8477b13f9f965f1a1e18524bf50881","impliedNodeFormat":1}],"fileIdsList":[[2]],"options":{"declaration":true,"module":99},"referencedMap":[[3,1]]}
//// [/home/src/workspaces/project/tsconfig.tsbuildinfo.readable.baseline.txt] *new* 
{
  "version": "FakeTSVersion",
  "fileNames": [
    "../../tslibs/TS/Lib/lib.d.ts",
    "./MessageablePerson.ts",
    "./main.ts"
  ],
  "fileInfos": [
    {
      "fileName": "../../tslibs/TS/Lib/lib.d.ts",
      "version": "575a4e15624573144926595079b1ec30f9c7853bab32f43c0b7db2acfdf038e2",
      "signature": "575a4e15624573144926595079b1ec30f9c7853bab32f43c0b7db2acfdf038e2",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "575a4e15624573144926595079b1ec30f9c7853bab32f43c0b7db2acfdf038e2",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./MessageablePerson.ts",
      "version": "ff666de4fdc53b5500de60a9b8c073c9327a9e9326417ef4861b8d2473c7457a",
      "signature": "6ec1f7bdc192ba06258caff3fa202fd577f8f354d676f548500eeb232155cbbe",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "ff666de4fdc53b5500de60a9b8c073c9327a9e9326417ef4861b8d2473c7457a",
        "signature": "6ec1f7bdc192ba06258caff3fa202fd577f8f354d676f548500eeb232155cbbe",
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./main.ts",
      "version": "36f0b00de3c707929bf1919e32e5b6053c8730bb00aa779bcdd1925414d68b8c",
      "signature": "8e609bb71c20b858c77f0e9f90bb1319db8477b13f9f965f1a1e18524bf50881",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "36f0b00de3c707929bf1919e32e5b6053c8730bb00aa779bcdd1925414d68b8c",
        "signature": "8e609bb71c20b858c77f0e9f90bb1319db8477b13f9f965f1a1e18524bf50881",
        "impliedNodeFormat": 1
      }
    }
  ],
  "fileIdsList": [
    [
      "./MessageablePerson.ts"
    ]
  ],
  "options": {
    "declaration": true,
    "module": 99
  },
  "referencedMap": {
    "./main.ts": [
      "./MessageablePerson.ts"
    ]
  },
  "size": 697
}

