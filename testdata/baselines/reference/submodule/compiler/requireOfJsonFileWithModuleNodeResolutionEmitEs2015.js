//// [tests/cases/compiler/requireOfJsonFileWithModuleNodeResolutionEmitEs2015.ts] ////

//// [file1.ts]
import * as b from './b.json';

//// [b.json]
{
    "a": true,
    "b": "hello"
}

out/b.json(1,1): error TS1005: '{' expected.
out/b.json(1,2): error TS1136: Property assignment expected.
out/b.json(4,2): error TS1012: Unexpected token.
out/b.json(5,1): error TS1005: '}' expected.


==== out/b.json (4 errors) ====
    ({
    ~
!!! error TS1005: '{' expected.
     ~
!!! error TS1136: Property assignment expected.
        "a": true,
        "b": "hello"
    })
     ~
!!! error TS1012: Unexpected token.
    
    
!!! error TS1005: '}' expected.
//// [out/file1.js]
export {};
