//// [tests/cases/compiler/commentsEnums.ts] ////

//// [commentsEnums.ts]
/** Enum of colors*/
enum Colors {
    /** Fancy name for 'blue'*/
    Cornflower /* blue */,
    /** Fancy name for 'pink'*/
    FancyPink
} // trailing comment
var x = Colors.Cornflower;
x = Colors.FancyPink;



//// [commentsEnums.js]
"use strict";
/** Enum of colors*/
var Colors;
(function (Colors) {
    /** Fancy name for 'blue'*/
    Colors[Colors["Cornflower"] = 0] = "Cornflower"; /* blue */
    /** Fancy name for 'pink'*/
    Colors[Colors["FancyPink"] = 1] = "FancyPink";
})(Colors || (Colors = {})); // trailing comment
var x = Colors.Cornflower;
x = Colors.FancyPink;


//// [commentsEnums.d.ts]
/** Enum of colors*/
enum Colors {
    /** Fancy name for 'blue'*/
    Cornflower = 0,
    /** Fancy name for 'pink'*/
    FancyPink = 1
}
var x: Colors;


//// [DtsFileErrors]


commentsEnums.d.ts(2,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== commentsEnums.d.ts (1 errors) ====
    /** Enum of colors*/
    enum Colors {
    ~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        /** Fancy name for 'blue'*/
        Cornflower = 0,
        /** Fancy name for 'pink'*/
        FancyPink = 1
    }
    var x: Colors;
    