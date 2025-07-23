currentDirectory::/home/src/workspaces/project
useCaseSensitiveFileNames::true
Input::
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
{ 
    "compilerOptions": { 
        "module": "esnext"
    }
}

tsgo --incremental
ExitStatus:: Success
Output::
//// [/home/src/workspaces/project/MessageablePerson.js] *new* 
const Messageable = () => {
    return class MessageableClass {
        message = 'hello';
    };
};
const wrapper = () => Messageable();
export {};

//// [/home/src/workspaces/project/main.js] *new* 
function logMessage(person) {
    console.log(person.message);
}
export {};

//// [/home/src/workspaces/project/tsconfig.tsbuildinfo] *new* 
{"version":"FakeTSVersion","fileNames":["../../tslibs/TS/Lib/lib.d.ts","./MessageablePerson.ts","./main.ts"],"fileInfos":[{"version":"4454fdb8db546b8967485a3a7254c948e6876fb850a20e51972933eaf60b5b21-/// \u003creference no-default-lib=\"true\"/\u003e\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array\u003cT\u003e { length: number; [n: number]: T; }\ninterface ReadonlyArray\u003cT\u003e {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };\ntype ReturnType\u003cT extends (...args: any) =\u003e any\u003e = T extends (...args: any) =\u003e infer R ? R : any;\ntype InstanceType\u003cT extends abstract new (...args: any) =\u003e any\u003e = T extends abstract new (...args: any) =\u003e infer R ? R : any;","affectsGlobalScope":true,"impliedNodeFormat":1},"3be6695caa91776ec738c01ffbc1250eb86f9bca0c22b02335b5e5c7c63bcbaa-const Messageable = () =\u003e {\n    return class MessageableClass {\n        public message = 'hello';\n    }\n};\nconst wrapper = () =\u003e Messageable();\ntype MessageablePerson = InstanceType\u003cReturnType\u003ctypeof wrapper\u003e\u003e;\nexport default MessageablePerson;","1fe1ef024191a0efa3b3b0ec73ee8e703a58d44f6d62caf49268591964dce1de-import MessageablePerson from './MessageablePerson.js';\nfunction logMessage( person: MessageablePerson ) {\n    console.log( person.message );\n}"],"fileIdsList":[[2]],"options":{"module":99},"referencedMap":[[3,1]]}
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
      "version": "4454fdb8db546b8967485a3a7254c948e6876fb850a20e51972933eaf60b5b21-/// \u003creference no-default-lib=\"true\"/\u003e\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array\u003cT\u003e { length: number; [n: number]: T; }\ninterface ReadonlyArray\u003cT\u003e {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };\ntype ReturnType\u003cT extends (...args: any) =\u003e any\u003e = T extends (...args: any) =\u003e infer R ? R : any;\ntype InstanceType\u003cT extends abstract new (...args: any) =\u003e any\u003e = T extends abstract new (...args: any) =\u003e infer R ? R : any;",
      "signature": "4454fdb8db546b8967485a3a7254c948e6876fb850a20e51972933eaf60b5b21-/// \u003creference no-default-lib=\"true\"/\u003e\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array\u003cT\u003e { length: number; [n: number]: T; }\ninterface ReadonlyArray\u003cT\u003e {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };\ntype ReturnType\u003cT extends (...args: any) =\u003e any\u003e = T extends (...args: any) =\u003e infer R ? R : any;\ntype InstanceType\u003cT extends abstract new (...args: any) =\u003e any\u003e = T extends abstract new (...args: any) =\u003e infer R ? R : any;",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "4454fdb8db546b8967485a3a7254c948e6876fb850a20e51972933eaf60b5b21-/// \u003creference no-default-lib=\"true\"/\u003e\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array\u003cT\u003e { length: number; [n: number]: T; }\ninterface ReadonlyArray\u003cT\u003e {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };\ntype ReturnType\u003cT extends (...args: any) =\u003e any\u003e = T extends (...args: any) =\u003e infer R ? R : any;\ntype InstanceType\u003cT extends abstract new (...args: any) =\u003e any\u003e = T extends abstract new (...args: any) =\u003e infer R ? R : any;",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./MessageablePerson.ts",
      "version": "3be6695caa91776ec738c01ffbc1250eb86f9bca0c22b02335b5e5c7c63bcbaa-const Messageable = () =\u003e {\n    return class MessageableClass {\n        public message = 'hello';\n    }\n};\nconst wrapper = () =\u003e Messageable();\ntype MessageablePerson = InstanceType\u003cReturnType\u003ctypeof wrapper\u003e\u003e;\nexport default MessageablePerson;",
      "signature": "3be6695caa91776ec738c01ffbc1250eb86f9bca0c22b02335b5e5c7c63bcbaa-const Messageable = () =\u003e {\n    return class MessageableClass {\n        public message = 'hello';\n    }\n};\nconst wrapper = () =\u003e Messageable();\ntype MessageablePerson = InstanceType\u003cReturnType\u003ctypeof wrapper\u003e\u003e;\nexport default MessageablePerson;",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "./main.ts",
      "version": "1fe1ef024191a0efa3b3b0ec73ee8e703a58d44f6d62caf49268591964dce1de-import MessageablePerson from './MessageablePerson.js';\nfunction logMessage( person: MessageablePerson ) {\n    console.log( person.message );\n}",
      "signature": "1fe1ef024191a0efa3b3b0ec73ee8e703a58d44f6d62caf49268591964dce1de-import MessageablePerson from './MessageablePerson.js';\nfunction logMessage( person: MessageablePerson ) {\n    console.log( person.message );\n}",
      "impliedNodeFormat": "CommonJS"
    }
  ],
  "fileIdsList": [
    [
      "./MessageablePerson.ts"
    ]
  ],
  "options": {
    "module": 99
  },
  "referencedMap": {
    "./main.ts": [
      "./MessageablePerson.ts"
    ]
  },
  "size": 1852
}

SemanticDiagnostics::
*refresh*    /home/src/tslibs/TS/Lib/lib.d.ts
*refresh*    /home/src/workspaces/project/MessageablePerson.ts
*refresh*    /home/src/workspaces/project/main.ts
Signatures::


Edit [0]:: no change

tsgo --incremental
ExitStatus:: Success
Output::

SemanticDiagnostics::
Signatures::


Edit [1]:: modify public to protected
//// [/home/src/workspaces/project/MessageablePerson.ts] *modified* 
const Messageable = () => {
    return class MessageableClass {
        protected message = 'hello';
    }
};
const wrapper = () => Messageable();
type MessageablePerson = InstanceType<ReturnType<typeof wrapper>>;
export default MessageablePerson;

tsgo --incremental
ExitStatus:: DiagnosticsPresent_OutputsGenerated
Output::
[96mmain.ts[0m:[93m3[0m:[93m25[0m - [91merror[0m[90m TS2445: [0mProperty 'message' is protected and only accessible within class 'MessageableClass' and its subclasses.

[7m3[0m     console.log( person.message );
[7m [0m [91m                        ~~~~~~~[0m


Found 1 error in main.ts[90m:3[0m

//// [/home/src/workspaces/project/MessageablePerson.js] *rewrite with same content*
//// [/home/src/workspaces/project/main.js] *rewrite with same content*
//// [/home/src/workspaces/project/tsconfig.tsbuildinfo] *modified* 
{"version":"FakeTSVersion","fileNames":["../../tslibs/TS/Lib/lib.d.ts","./MessageablePerson.ts","./main.ts"],"fileInfos":[{"version":"4454fdb8db546b8967485a3a7254c948e6876fb850a20e51972933eaf60b5b21-/// \u003creference no-default-lib=\"true\"/\u003e\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array\u003cT\u003e { length: number; [n: number]: T; }\ninterface ReadonlyArray\u003cT\u003e {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };\ntype ReturnType\u003cT extends (...args: any) =\u003e any\u003e = T extends (...args: any) =\u003e infer R ? R : any;\ntype InstanceType\u003cT extends abstract new (...args: any) =\u003e any\u003e = T extends abstract new (...args: any) =\u003e infer R ? R : any;","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"7605d26cb18fccfe18cf0d5b1771a538d665f6b2bb331c9617ebdf38d7c93b29-const Messageable = () =\u003e {\n    return class MessageableClass {\n        protected message = 'hello';\n    }\n};\nconst wrapper = () =\u003e Messageable();\ntype MessageablePerson = InstanceType\u003cReturnType\u003ctypeof wrapper\u003e\u003e;\nexport default MessageablePerson;","signature":"34f86a1082929a2d2c6b784bde116fadfb61a9ee55f5141d4906ef4ce16a89c9-declare const wrapper: () =\u003e {\n    new (): {\n        message: string;\n    };\n};\ntype MessageablePerson = InstanceType\u003cReturnType\u003ctypeof wrapper\u003e\u003e;\nexport default MessageablePerson;\n\n(116,7): error4094: Property 'message' of exported anonymous class type may not be private or protected.\n(116,7): error9027: Add a type annotation to the variable wrapper.","impliedNodeFormat":1},{"version":"1fe1ef024191a0efa3b3b0ec73ee8e703a58d44f6d62caf49268591964dce1de-import MessageablePerson from './MessageablePerson.js';\nfunction logMessage( person: MessageablePerson ) {\n    console.log( person.message );\n}","signature":"8e609bb71c20b858c77f0e9f90bb1319db8477b13f9f965f1a1e18524bf50881-export {};\n","impliedNodeFormat":1}],"fileIdsList":[[2]],"options":{"module":99},"referencedMap":[[3,1]],"semanticDiagnosticsPerFile":[[3,[{"pos":131,"end":138,"code":2445,"category":1,"message":"Property 'message' is protected and only accessible within class 'MessageableClass' and its subclasses."}]]]}
//// [/home/src/workspaces/project/tsconfig.tsbuildinfo.readable.baseline.txt] *modified* 
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
      "version": "4454fdb8db546b8967485a3a7254c948e6876fb850a20e51972933eaf60b5b21-/// \u003creference no-default-lib=\"true\"/\u003e\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array\u003cT\u003e { length: number; [n: number]: T; }\ninterface ReadonlyArray\u003cT\u003e {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };\ntype ReturnType\u003cT extends (...args: any) =\u003e any\u003e = T extends (...args: any) =\u003e infer R ? R : any;\ntype InstanceType\u003cT extends abstract new (...args: any) =\u003e any\u003e = T extends abstract new (...args: any) =\u003e infer R ? R : any;",
      "signature": "4454fdb8db546b8967485a3a7254c948e6876fb850a20e51972933eaf60b5b21-/// \u003creference no-default-lib=\"true\"/\u003e\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array\u003cT\u003e { length: number; [n: number]: T; }\ninterface ReadonlyArray\u003cT\u003e {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };\ntype ReturnType\u003cT extends (...args: any) =\u003e any\u003e = T extends (...args: any) =\u003e infer R ? R : any;\ntype InstanceType\u003cT extends abstract new (...args: any) =\u003e any\u003e = T extends abstract new (...args: any) =\u003e infer R ? R : any;",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "4454fdb8db546b8967485a3a7254c948e6876fb850a20e51972933eaf60b5b21-/// \u003creference no-default-lib=\"true\"/\u003e\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array\u003cT\u003e { length: number; [n: number]: T; }\ninterface ReadonlyArray\u003cT\u003e {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };\ntype ReturnType\u003cT extends (...args: any) =\u003e any\u003e = T extends (...args: any) =\u003e infer R ? R : any;\ntype InstanceType\u003cT extends abstract new (...args: any) =\u003e any\u003e = T extends abstract new (...args: any) =\u003e infer R ? R : any;",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./MessageablePerson.ts",
      "version": "7605d26cb18fccfe18cf0d5b1771a538d665f6b2bb331c9617ebdf38d7c93b29-const Messageable = () =\u003e {\n    return class MessageableClass {\n        protected message = 'hello';\n    }\n};\nconst wrapper = () =\u003e Messageable();\ntype MessageablePerson = InstanceType\u003cReturnType\u003ctypeof wrapper\u003e\u003e;\nexport default MessageablePerson;",
      "signature": "34f86a1082929a2d2c6b784bde116fadfb61a9ee55f5141d4906ef4ce16a89c9-declare const wrapper: () =\u003e {\n    new (): {\n        message: string;\n    };\n};\ntype MessageablePerson = InstanceType\u003cReturnType\u003ctypeof wrapper\u003e\u003e;\nexport default MessageablePerson;\n\n(116,7): error4094: Property 'message' of exported anonymous class type may not be private or protected.\n(116,7): error9027: Add a type annotation to the variable wrapper.",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "7605d26cb18fccfe18cf0d5b1771a538d665f6b2bb331c9617ebdf38d7c93b29-const Messageable = () =\u003e {\n    return class MessageableClass {\n        protected message = 'hello';\n    }\n};\nconst wrapper = () =\u003e Messageable();\ntype MessageablePerson = InstanceType\u003cReturnType\u003ctypeof wrapper\u003e\u003e;\nexport default MessageablePerson;",
        "signature": "34f86a1082929a2d2c6b784bde116fadfb61a9ee55f5141d4906ef4ce16a89c9-declare const wrapper: () =\u003e {\n    new (): {\n        message: string;\n    };\n};\ntype MessageablePerson = InstanceType\u003cReturnType\u003ctypeof wrapper\u003e\u003e;\nexport default MessageablePerson;\n\n(116,7): error4094: Property 'message' of exported anonymous class type may not be private or protected.\n(116,7): error9027: Add a type annotation to the variable wrapper.",
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./main.ts",
      "version": "1fe1ef024191a0efa3b3b0ec73ee8e703a58d44f6d62caf49268591964dce1de-import MessageablePerson from './MessageablePerson.js';\nfunction logMessage( person: MessageablePerson ) {\n    console.log( person.message );\n}",
      "signature": "8e609bb71c20b858c77f0e9f90bb1319db8477b13f9f965f1a1e18524bf50881-export {};\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "1fe1ef024191a0efa3b3b0ec73ee8e703a58d44f6d62caf49268591964dce1de-import MessageablePerson from './MessageablePerson.js';\nfunction logMessage( person: MessageablePerson ) {\n    console.log( person.message );\n}",
        "signature": "8e609bb71c20b858c77f0e9f90bb1319db8477b13f9f965f1a1e18524bf50881-export {};\n",
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
    "module": 99
  },
  "referencedMap": {
    "./main.ts": [
      "./MessageablePerson.ts"
    ]
  },
  "semanticDiagnosticsPerFile": [
    [
      "./main.ts",
      [
        {
          "pos": 131,
          "end": 138,
          "code": 2445,
          "category": 1,
          "message": "Property 'message' is protected and only accessible within class 'MessageableClass' and its subclasses."
        }
      ]
    ]
  ],
  "size": 2682
}

SemanticDiagnostics::
*refresh*    /home/src/workspaces/project/MessageablePerson.ts
*refresh*    /home/src/workspaces/project/main.ts
Signatures::
(computed .d.ts) /home/src/workspaces/project/MessageablePerson.ts
(computed .d.ts) /home/src/workspaces/project/main.ts


Edit [2]:: no change

tsgo --incremental
ExitStatus:: DiagnosticsPresent_OutputsGenerated
Output::
[96mmain.ts[0m:[93m3[0m:[93m25[0m - [91merror[0m[90m TS2445: [0mProperty 'message' is protected and only accessible within class 'MessageableClass' and its subclasses.

[7m3[0m     console.log( person.message );
[7m [0m [91m                        ~~~~~~~[0m


Found 1 error in main.ts[90m:3[0m


SemanticDiagnostics::
Signatures::


Edit [3]:: modify protected to public
//// [/home/src/workspaces/project/MessageablePerson.ts] *modified* 
const Messageable = () => {
    return class MessageableClass {
        public message = 'hello';
    }
};
const wrapper = () => Messageable();
type MessageablePerson = InstanceType<ReturnType<typeof wrapper>>;
export default MessageablePerson;

tsgo --incremental
ExitStatus:: Success
Output::
//// [/home/src/workspaces/project/MessageablePerson.js] *rewrite with same content*
//// [/home/src/workspaces/project/main.js] *rewrite with same content*
//// [/home/src/workspaces/project/tsconfig.tsbuildinfo] *modified* 
{"version":"FakeTSVersion","fileNames":["../../tslibs/TS/Lib/lib.d.ts","./MessageablePerson.ts","./main.ts"],"fileInfos":[{"version":"4454fdb8db546b8967485a3a7254c948e6876fb850a20e51972933eaf60b5b21-/// \u003creference no-default-lib=\"true\"/\u003e\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array\u003cT\u003e { length: number; [n: number]: T; }\ninterface ReadonlyArray\u003cT\u003e {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };\ntype ReturnType\u003cT extends (...args: any) =\u003e any\u003e = T extends (...args: any) =\u003e infer R ? R : any;\ntype InstanceType\u003cT extends abstract new (...args: any) =\u003e any\u003e = T extends abstract new (...args: any) =\u003e infer R ? R : any;","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"3be6695caa91776ec738c01ffbc1250eb86f9bca0c22b02335b5e5c7c63bcbaa-const Messageable = () =\u003e {\n    return class MessageableClass {\n        public message = 'hello';\n    }\n};\nconst wrapper = () =\u003e Messageable();\ntype MessageablePerson = InstanceType\u003cReturnType\u003ctypeof wrapper\u003e\u003e;\nexport default MessageablePerson;","signature":"6ec1f7bdc192ba06258caff3fa202fd577f8f354d676f548500eeb232155cbbe-declare const wrapper: () =\u003e {\n    new (): {\n        message: string;\n    };\n};\ntype MessageablePerson = InstanceType\u003cReturnType\u003ctypeof wrapper\u003e\u003e;\nexport default MessageablePerson;\n","impliedNodeFormat":1},{"version":"1fe1ef024191a0efa3b3b0ec73ee8e703a58d44f6d62caf49268591964dce1de-import MessageablePerson from './MessageablePerson.js';\nfunction logMessage( person: MessageablePerson ) {\n    console.log( person.message );\n}","signature":"8e609bb71c20b858c77f0e9f90bb1319db8477b13f9f965f1a1e18524bf50881-export {};\n","impliedNodeFormat":1}],"fileIdsList":[[2]],"options":{"module":99},"referencedMap":[[3,1]]}
//// [/home/src/workspaces/project/tsconfig.tsbuildinfo.readable.baseline.txt] *modified* 
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
      "version": "4454fdb8db546b8967485a3a7254c948e6876fb850a20e51972933eaf60b5b21-/// \u003creference no-default-lib=\"true\"/\u003e\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array\u003cT\u003e { length: number; [n: number]: T; }\ninterface ReadonlyArray\u003cT\u003e {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };\ntype ReturnType\u003cT extends (...args: any) =\u003e any\u003e = T extends (...args: any) =\u003e infer R ? R : any;\ntype InstanceType\u003cT extends abstract new (...args: any) =\u003e any\u003e = T extends abstract new (...args: any) =\u003e infer R ? R : any;",
      "signature": "4454fdb8db546b8967485a3a7254c948e6876fb850a20e51972933eaf60b5b21-/// \u003creference no-default-lib=\"true\"/\u003e\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array\u003cT\u003e { length: number; [n: number]: T; }\ninterface ReadonlyArray\u003cT\u003e {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };\ntype ReturnType\u003cT extends (...args: any) =\u003e any\u003e = T extends (...args: any) =\u003e infer R ? R : any;\ntype InstanceType\u003cT extends abstract new (...args: any) =\u003e any\u003e = T extends abstract new (...args: any) =\u003e infer R ? R : any;",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "4454fdb8db546b8967485a3a7254c948e6876fb850a20e51972933eaf60b5b21-/// \u003creference no-default-lib=\"true\"/\u003e\ninterface Boolean {}\ninterface Function {}\ninterface CallableFunction {}\ninterface NewableFunction {}\ninterface IArguments {}\ninterface Number { toExponential: any; }\ninterface Object {}\ninterface RegExp {}\ninterface String { charAt: any; }\ninterface Array\u003cT\u003e { length: number; [n: number]: T; }\ninterface ReadonlyArray\u003cT\u003e {}\ninterface SymbolConstructor {\n    (desc?: string | number): symbol;\n    for(name: string): symbol;\n    readonly toStringTag: symbol;\n}\ndeclare var Symbol: SymbolConstructor;\ninterface Symbol {\n    readonly [Symbol.toStringTag]: string;\n}\ndeclare const console: { log(msg: any): void; };\ntype ReturnType\u003cT extends (...args: any) =\u003e any\u003e = T extends (...args: any) =\u003e infer R ? R : any;\ntype InstanceType\u003cT extends abstract new (...args: any) =\u003e any\u003e = T extends abstract new (...args: any) =\u003e infer R ? R : any;",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./MessageablePerson.ts",
      "version": "3be6695caa91776ec738c01ffbc1250eb86f9bca0c22b02335b5e5c7c63bcbaa-const Messageable = () =\u003e {\n    return class MessageableClass {\n        public message = 'hello';\n    }\n};\nconst wrapper = () =\u003e Messageable();\ntype MessageablePerson = InstanceType\u003cReturnType\u003ctypeof wrapper\u003e\u003e;\nexport default MessageablePerson;",
      "signature": "6ec1f7bdc192ba06258caff3fa202fd577f8f354d676f548500eeb232155cbbe-declare const wrapper: () =\u003e {\n    new (): {\n        message: string;\n    };\n};\ntype MessageablePerson = InstanceType\u003cReturnType\u003ctypeof wrapper\u003e\u003e;\nexport default MessageablePerson;\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "3be6695caa91776ec738c01ffbc1250eb86f9bca0c22b02335b5e5c7c63bcbaa-const Messageable = () =\u003e {\n    return class MessageableClass {\n        public message = 'hello';\n    }\n};\nconst wrapper = () =\u003e Messageable();\ntype MessageablePerson = InstanceType\u003cReturnType\u003ctypeof wrapper\u003e\u003e;\nexport default MessageablePerson;",
        "signature": "6ec1f7bdc192ba06258caff3fa202fd577f8f354d676f548500eeb232155cbbe-declare const wrapper: () =\u003e {\n    new (): {\n        message: string;\n    };\n};\ntype MessageablePerson = InstanceType\u003cReturnType\u003ctypeof wrapper\u003e\u003e;\nexport default MessageablePerson;\n",
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "./main.ts",
      "version": "1fe1ef024191a0efa3b3b0ec73ee8e703a58d44f6d62caf49268591964dce1de-import MessageablePerson from './MessageablePerson.js';\nfunction logMessage( person: MessageablePerson ) {\n    console.log( person.message );\n}",
      "signature": "8e609bb71c20b858c77f0e9f90bb1319db8477b13f9f965f1a1e18524bf50881-export {};\n",
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "1fe1ef024191a0efa3b3b0ec73ee8e703a58d44f6d62caf49268591964dce1de-import MessageablePerson from './MessageablePerson.js';\nfunction logMessage( person: MessageablePerson ) {\n    console.log( person.message );\n}",
        "signature": "8e609bb71c20b858c77f0e9f90bb1319db8477b13f9f965f1a1e18524bf50881-export {};\n",
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
    "module": 99
  },
  "referencedMap": {
    "./main.ts": [
      "./MessageablePerson.ts"
    ]
  },
  "size": 2305
}

SemanticDiagnostics::
*refresh*    /home/src/workspaces/project/MessageablePerson.ts
*refresh*    /home/src/workspaces/project/main.ts
Signatures::
(computed .d.ts) /home/src/workspaces/project/MessageablePerson.ts
(computed .d.ts) /home/src/workspaces/project/main.ts


Edit [4]:: no change

tsgo --incremental
ExitStatus:: Success
Output::

SemanticDiagnostics::
Signatures::
