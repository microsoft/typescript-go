//// [tests/cases/compiler/forLoopInitializerAlwaysThrowsInTryCatch.ts] ////

//// [forLoopInitializerAlwaysThrowsInTryCatch.ts]
// Test 1: The original crash case
function test1() {
    try {
        for ((() => { throw new Error(); })(); ;) {
        }
    } catch (e) {
    }
}

// Test 2: Unlabelled break and continue inside the unreachable body
function test2() {
    try {
        for ((() => { throw new Error(); })(); ;) {
            if (Math.random()) {
                break;
            } else {
                continue;
            }
        }
    } catch (e) {
    }
}

// Test 3: Nested loops where inner loop has unreachable body with break/continue
// Verifies that break/continue don't leak and pollute the outer loop's CFG
function test3() {
    while (true) {
        try {
            for ((() => { throw new Error(); })(); ;) {
                if (Math.random()) {
                    break;
                } else {
                    continue;
                }
            }
        } catch (e) {
        }
    }
}

// Test 4: Labeled break and continue inside the unreachable body
function test4() {
    try {
        label1: for ((() => { throw new Error(); })(); ;) {
            if (Math.random()) {
                break label1;
            } else {
                continue label1;
            }
        }
    } catch (e) {
    }
}


//// [forLoopInitializerAlwaysThrowsInTryCatch.js]
"use strict";
// Test 1: The original crash case
function test1() {
    try {
        for ((() => { throw new Error(); })();;) {
        }
    }
    catch (e) {
    }
}
// Test 2: Unlabelled break and continue inside the unreachable body
function test2() {
    try {
        for ((() => { throw new Error(); })();;) {
            if (Math.random()) {
                break;
            }
            else {
                continue;
            }
        }
    }
    catch (e) {
    }
}
// Test 3: Nested loops where inner loop has unreachable body with break/continue
// Verifies that break/continue don't leak and pollute the outer loop's CFG
function test3() {
    while (true) {
        try {
            for ((() => { throw new Error(); })();;) {
                if (Math.random()) {
                    break;
                }
                else {
                    continue;
                }
            }
        }
        catch (e) {
        }
    }
}
// Test 4: Labeled break and continue inside the unreachable body
function test4() {
    try {
        label1: for ((() => { throw new Error(); })();;) {
            if (Math.random()) {
                break label1;
            }
            else {
                continue label1;
            }
        }
    }
    catch (e) {
    }
}
