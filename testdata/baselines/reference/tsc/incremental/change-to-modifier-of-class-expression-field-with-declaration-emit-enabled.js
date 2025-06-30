
currentDirectory::/home/src/workspaces/project
useCaseSensitiveFileNames::true
Input::--incremental
//// [/home/src/workspaces/project/MessageablePerson.ts] new file

                        const Messageable = () => {
                            return class MessageableClass {
                                public message = 'hello';
                            }
                        };
                        const wrapper = () => Messageable();
                        type MessageablePerson = InstanceType<ReturnType<typeof wrapper>>;
                        export default MessageablePerson;
//// [/home/src/workspaces/project/main.ts] new file

                        import MessageablePerson from './MessageablePerson.js';
                        function logMessage( person: MessageablePerson ) {
                            console.log( person.message );
                        }
//// [/home/src/workspaces/project/tsconfig.json] new file
{ "compilerOptions": { "module": "esnext", "declaration": true  } }

ExitStatus:: 0

CompilerOptions::{
    "incremental": true
}
Output::
//// [/home/src/workspaces/project/MessageablePerson.d.ts] new file
declare const wrapper: () => {
    new (): {
        message: string;
    };
};
type MessageablePerson = InstanceType<ReturnType<typeof wrapper>>;
export default MessageablePerson;

//// [/home/src/workspaces/project/MessageablePerson.js] new file
const Messageable = () => {
    return class MessageableClass {
        message = 'hello';
    };
};
const wrapper = () => Messageable();
export {};

//// [/home/src/workspaces/project/MessageablePerson.ts] no change
//// [/home/src/workspaces/project/main.d.ts] new file
export {};

//// [/home/src/workspaces/project/main.js] new file
function logMessage(person) {
    console.log(person.message);
}
export {};

//// [/home/src/workspaces/project/main.ts] no change
//// [/home/src/workspaces/project/tsconfig.json] no change
//// [/home/src/workspaces/project/tsconfig.tsbuildinfo] new file
{"Version":"FakeTSVersion","fileNames":["bundled:///libs/lib.d.ts","bundled:///libs/lib.es5.d.ts","bundled:///libs/lib.dom.d.ts","bundled:///libs/lib.webworker.importscripts.d.ts","bundled:///libs/lib.scripthost.d.ts","bundled:///libs/lib.decorators.d.ts","bundled:///libs/lib.decorators.legacy.d.ts","./MessageablePerson.ts","./main.ts"],"fileInfos":["a7297ff837fcdf174a9524925966429eb8e5feecc2cc55cc06574e6b092c1eaa",{"version":"69684132aeb9b5642cbcd9e22dff7818ff0ee1aa831728af0ecf97d3364d5546","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"092c2bfe125ce69dbb1223c85d68d4d2397d7d8411867b5cc03cec902c233763","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"80e18897e5884b6723488d4f5652167e7bb5024f946743134ecc4aa4ee731f89","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"cd034f499c6cdca722b60c04b5b1b78e058487a7085a8e0d6fb50809947ee573","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"8e7f8264d0fb4c5339605a15daadb037bf238c10b654bb3eee14208f860a32ea","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"782dec38049b92d4e85c1585fbea5474a219c6984a35b004963b00beb1aab538","affectsGlobalScope":true,"impliedNodeFormat":1},{"version":"ff666de4fdc53b5500de60a9b8c073c9327a9e9326417ef4861b8d2473c7457a","signature":"6ec1f7bdc192ba06258caff3fa202fd577f8f354d676f548500eeb232155cbbe","impliedNodeFormat":1},{"version":"36f0b00de3c707929bf1919e32e5b6053c8730bb00aa779bcdd1925414d68b8c","signature":"8e609bb71c20b858c77f0e9f90bb1319db8477b13f9f965f1a1e18524bf50881","impliedNodeFormat":1}],"fileIdsList":[[8]],"options":{"declaration":true,"module":99},"referencedMap":[[9,1]]}
//// [/home/src/workspaces/project/tsconfig.tsbuildinfo.readable.baseline.txt] new file
{
  "Version": "FakeTSVersion",
  "fileNames": [
    "bundled:///libs/lib.d.ts",
    "bundled:///libs/lib.es5.d.ts",
    "bundled:///libs/lib.dom.d.ts",
    "bundled:///libs/lib.webworker.importscripts.d.ts",
    "bundled:///libs/lib.scripthost.d.ts",
    "bundled:///libs/lib.decorators.d.ts",
    "bundled:///libs/lib.decorators.legacy.d.ts",
    "./MessageablePerson.ts",
    "./main.ts"
  ],
  "fileInfos": [
    {
      "fileName": "bundled:///libs/lib.d.ts",
      "version": "a7297ff837fcdf174a9524925966429eb8e5feecc2cc55cc06574e6b092c1eaa",
      "signature": "a7297ff837fcdf174a9524925966429eb8e5feecc2cc55cc06574e6b092c1eaa",
      "impliedNodeFormat": "CommonJS"
    },
    {
      "fileName": "bundled:///libs/lib.es5.d.ts",
      "version": "69684132aeb9b5642cbcd9e22dff7818ff0ee1aa831728af0ecf97d3364d5546",
      "signature": "69684132aeb9b5642cbcd9e22dff7818ff0ee1aa831728af0ecf97d3364d5546",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "69684132aeb9b5642cbcd9e22dff7818ff0ee1aa831728af0ecf97d3364d5546",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "bundled:///libs/lib.dom.d.ts",
      "version": "092c2bfe125ce69dbb1223c85d68d4d2397d7d8411867b5cc03cec902c233763",
      "signature": "092c2bfe125ce69dbb1223c85d68d4d2397d7d8411867b5cc03cec902c233763",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "092c2bfe125ce69dbb1223c85d68d4d2397d7d8411867b5cc03cec902c233763",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "bundled:///libs/lib.webworker.importscripts.d.ts",
      "version": "80e18897e5884b6723488d4f5652167e7bb5024f946743134ecc4aa4ee731f89",
      "signature": "80e18897e5884b6723488d4f5652167e7bb5024f946743134ecc4aa4ee731f89",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "80e18897e5884b6723488d4f5652167e7bb5024f946743134ecc4aa4ee731f89",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "bundled:///libs/lib.scripthost.d.ts",
      "version": "cd034f499c6cdca722b60c04b5b1b78e058487a7085a8e0d6fb50809947ee573",
      "signature": "cd034f499c6cdca722b60c04b5b1b78e058487a7085a8e0d6fb50809947ee573",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "cd034f499c6cdca722b60c04b5b1b78e058487a7085a8e0d6fb50809947ee573",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "bundled:///libs/lib.decorators.d.ts",
      "version": "8e7f8264d0fb4c5339605a15daadb037bf238c10b654bb3eee14208f860a32ea",
      "signature": "8e7f8264d0fb4c5339605a15daadb037bf238c10b654bb3eee14208f860a32ea",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "8e7f8264d0fb4c5339605a15daadb037bf238c10b654bb3eee14208f860a32ea",
        "affectsGlobalScope": true,
        "impliedNodeFormat": 1
      }
    },
    {
      "fileName": "bundled:///libs/lib.decorators.legacy.d.ts",
      "version": "782dec38049b92d4e85c1585fbea5474a219c6984a35b004963b00beb1aab538",
      "signature": "782dec38049b92d4e85c1585fbea5474a219c6984a35b004963b00beb1aab538",
      "affectsGlobalScope": true,
      "impliedNodeFormat": "CommonJS",
      "original": {
        "version": "782dec38049b92d4e85c1585fbea5474a219c6984a35b004963b00beb1aab538",
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
  "size": 1629
}

