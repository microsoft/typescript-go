//// [tests/cases/compiler/emitBundleWithShebang2.ts] ////

//// [test.ts]
#!/usr/bin/env gjs
class Doo {}
class Scooby extends Doo {}

//// [test2.ts]
#!/usr/bin/env js
class Dood {}
class Scoobyd extends Dood {}

//// [test2.js]
class Dood {
}
class Scoobyd extends Dood {
}
//// [test.js]
class Doo {
}
class Scooby extends Doo {
}
