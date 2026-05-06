//// [tests/cases/conformance/types/stringLiteral/stringLiteralTypesAndTuples01.ts] ////

//// [stringLiteralTypesAndTuples01.ts]
// Should all be strings.
let [hello, brave, newish, world] = ["Hello", "Brave", "New", "World"];

type RexOrRaptor = "t-rex" | "raptor"
let [im, a, dinosaur]: ["I'm", "a", RexOrRaptor] = ['I\'m', 'a', 't-rex'];

rawr(dinosaur);

function rawr(dino: RexOrRaptor) {
    if (dino === "t-rex") {
        return "ROAAAAR!";
    }
    if (dino === "raptor") {
        return "yip yip!";
    }

    throw "Unexpected " + dino;
}

//// [stringLiteralTypesAndTuples01.js]
"use strict";
// Should all be strings.
let [hello, brave, newish, world] = ["Hello", "Brave", "New", "World"];
let [im, a, dinosaur] = ['I\'m', 'a', 't-rex'];
rawr(dinosaur);
function rawr(dino) {
    if (dino === "t-rex") {
        return "ROAAAAR!";
    }
    if (dino === "raptor") {
        return "yip yip!";
    }
    throw "Unexpected " + dino;
}


//// [stringLiteralTypesAndTuples01.d.ts]
let hello: string, brave: string, newish: string, world: string;
type RexOrRaptor = "t-rex" | "raptor";
let im: "I'm", a: "a", dinosaur: RexOrRaptor;
function rawr(dino: RexOrRaptor): "ROAAAAR!" | "yip yip!";


//// [DtsFileErrors]


stringLiteralTypesAndTuples01.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== stringLiteralTypesAndTuples01.d.ts (1 errors) ====
    let hello: string, brave: string, newish: string, world: string;
    ~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    type RexOrRaptor = "t-rex" | "raptor";
    let im: "I'm", a: "a", dinosaur: RexOrRaptor;
    function rawr(dino: RexOrRaptor): "ROAAAAR!" | "yip yip!";
    