// @incremental: true
// @target: es2015
// @lib: es2015
// @tsBuildInfoFile: /a.tsbuildinfo

// Regression test: when a class extends a mixin that returns an intersection,
// the class's anonymous constructor type inherits static accessor properties from the
// intersection. These synthetic properties can have a nil parent symbol when the
// constituent accessor declarations differ. Using --incremental without --declaration,
// the incremental build must compute DTS signatures for changed files via forced
// declaration emit, which hits the same crash path.

// @filename: /a.tsbuildinfo
{
  "version": "TEST",
  "root": [14],
  "fileNames": [
    "lib.es5.d.ts",
    "lib.es2015.d.ts",
    "lib.es2015.core.d.ts",
    "lib.es2015.collection.d.ts",
    "lib.es2015.generator.d.ts",
    "lib.es2015.iterable.d.ts",
    "lib.es2015.promise.d.ts",
    "lib.es2015.proxy.d.ts",
    "lib.es2015.reflect.d.ts",
    "lib.es2015.symbol.d.ts",
    "lib.es2015.symbol.wellknown.d.ts",
    "lib.decorators.d.ts",
    "lib.decorators.legacy.d.ts",
    "./a.ts"
  ],
  "fileInfos": [
    {"version": "a1aa1a5e065d48ef5c7bb99e38412f96", "affectsGlobalScope": true, "impliedNodeFormat": 1},
    "d4306fb2e47f74835e8674ffac07d76f",
    {"version": "01ac052ec4a79e87229f90466a9645f8", "affectsGlobalScope": true, "impliedNodeFormat": 1},
    {"version": "edba5df642941aa062a62f6328c6df3d", "affectsGlobalScope": true, "impliedNodeFormat": 1},
    {"version": "6344b55f26a4e81d9608777dbfb877dd", "affectsGlobalScope": true, "impliedNodeFormat": 1},
    {"version": "3c0ed28e53d3695b363e256ec1c023fd", "affectsGlobalScope": true, "impliedNodeFormat": 1},
    {"version": "4c2761daba7f17141c25baa0821ac5da", "affectsGlobalScope": true, "impliedNodeFormat": 1},
    {"version": "b87656acabd63e69379ff6ffcfe52fc7", "affectsGlobalScope": true, "impliedNodeFormat": 1},
    {"version": "597469522da047a5af5222cc6989f405", "affectsGlobalScope": true, "impliedNodeFormat": 1},
    {"version": "bb3a710cbcda0533bb127712927cbe37", "affectsGlobalScope": true, "impliedNodeFormat": 1},
    {"version": "55d97a8c6fbf34a30450a7b1e5f7a298", "affectsGlobalScope": true, "impliedNodeFormat": 1},
    {"version": "f64453cbf9671f28158677fa5c43967a", "affectsGlobalScope": true, "impliedNodeFormat": 1},
    {"version": "33f317af5428801f944a478d2c1e38e5", "affectsGlobalScope": true, "impliedNodeFormat": 1},
    "stale-version-hash-to-force-recheck"
  ],
  "options": {
    "newLine": 1,
    "noErrorTruncation": true,
    "skipDefaultLibCheck": true,
    "target": 2,
    "tsBuildInfoFile": "./a.tsbuildinfo"
  },
  "changeFileSet": [14]
}
// @filename: /a.ts
function mix<T extends new (...args: any[]) => any, U extends new (...args: any[]) => any>(
    base1: T, base2: U
): T & U {
    return null as any;
}

class A {
    static get shared(): number { return 1; }
    static set shared(v: number) { }
    x: string = "";
}

class B {
    static get shared(): number { return 2; }
    static set shared(v: number) { }
    y: number = 0;
}

function make() {
    class C extends mix(A, B) {
        z: boolean = true;
    }
    return C;
}

export const MixedClass = make();
