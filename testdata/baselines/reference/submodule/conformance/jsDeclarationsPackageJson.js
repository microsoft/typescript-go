//// [tests/cases/conformance/jsdoc/declarations/jsDeclarationsPackageJson.ts] ////

//// [index.js]
const j = require("./package.json");
module.exports = j;
//// [package.json]
{
    "name": "pkg",
    "version": "0.1.0",
    "description": "A package",
    "main": "./dist/index.js",
    "bin": {
      "cli": "./bin/cli.js",
    },
    "engines": {
      "node": ">=0"
    },
    "scripts": {
      "scriptname": "run && run again",
    },
    "devDependencies": {
      "@ns/dep": "0.1.2",
    },
    "dependencies": {
      "dep": "1.2.3",
    },
    "repository": "microsoft/TypeScript",
    "keywords": [
      "kw"
    ],
    "author": "Auth",
    "license": "See Licensce",
    "homepage": "https://site",
    "config": {
      "o": ["a"]
    }
}
  

out/package.json(1,1): error TS1005: '{' expected.
out/package.json(1,2): error TS1136: Property assignment expected.
out/package.json(31,2): error TS1012: Unexpected token.
out/package.json(32,1): error TS1005: '}' expected.


==== out/package.json (4 errors) ====
    ({
    ~
!!! error TS1005: '{' expected.
     ~
!!! error TS1136: Property assignment expected.
        "name": "pkg",
        "version": "0.1.0",
        "description": "A package",
        "main": "./dist/index.js",
        "bin": {
            "cli": "./bin/cli.js"
        },
        "engines": {
            "node": ">=0"
        },
        "scripts": {
            "scriptname": "run && run again"
        },
        "devDependencies": {
            "@ns/dep": "0.1.2"
        },
        "dependencies": {
            "dep": "1.2.3"
        },
        "repository": "microsoft/TypeScript",
        "keywords": [
            "kw"
        ],
        "author": "Auth",
        "license": "See Licensce",
        "homepage": "https://site",
        "config": {
            "o": ["a"]
        }
    })
     ~
!!! error TS1012: Unexpected token.
    
    
!!! error TS1005: '}' expected.
//// [index.js]
const j = require("./package.json");
export = j;
module.exports = j;


//// [index.d.ts]
export = j;


//// [DtsFileErrors]


out/index.d.ts(1,10): error TS2304: Cannot find name 'j'.


==== out/index.d.ts (1 errors) ====
    export = j;
             ~
!!! error TS2304: Cannot find name 'j'.
    
==== package.json (0 errors) ====
    {
        "name": "pkg",
        "version": "0.1.0",
        "description": "A package",
        "main": "./dist/index.js",
        "bin": {
          "cli": "./bin/cli.js",
        },
        "engines": {
          "node": ">=0"
        },
        "scripts": {
          "scriptname": "run && run again",
        },
        "devDependencies": {
          "@ns/dep": "0.1.2",
        },
        "dependencies": {
          "dep": "1.2.3",
        },
        "repository": "microsoft/TypeScript",
        "keywords": [
          "kw"
        ],
        "author": "Auth",
        "license": "See Licensce",
        "homepage": "https://site",
        "config": {
          "o": ["a"]
        }
    }
      