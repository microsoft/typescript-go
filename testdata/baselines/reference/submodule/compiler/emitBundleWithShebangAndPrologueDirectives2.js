//// [tests/cases/compiler/emitBundleWithShebangAndPrologueDirectives2.ts] ////

//// [test.ts]
#!/usr/bin/env gjs
"use strict"
class Doo {}
class Scooby extends Doo {}

//// [test1.ts]
#!/usr/bin/env gjs
"use strict"
"Another prologue"
class Dood {}
class Scoobyd extends Dood {}

//// [test.js]
"use strict";
class Doo {
}
class Scooby extends Doo {
}
//// [test1.js]
"use strict";
"Another prologue";
class Dood {
}
class Scoobyd extends Dood {
}
