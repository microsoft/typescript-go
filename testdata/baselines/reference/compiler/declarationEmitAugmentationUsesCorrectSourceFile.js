//// [tests/cases/compiler/declarationEmitAugmentationUsesCorrectSourceFile.ts] ////

//// [index.d.ts]
// A bunch of random text to move the positions forward
// A bunch of random text to move the positions forward
// A bunch of random text to move the positions forward
// A bunch of random text to move the positions forward
// A bunch of random text to move the positions forward
// A bunch of random text to move the positions forward
// A bunch of random text to move the positions forward
// A bunch of random text to move the positions forward
// A bunch of random text to move the positions forward

type ShouldJustBeAny = [any][0];

declare namespace knex {
  export { Knex };
}

declare namespace Knex {
  interface Interface {
    method(): ShouldJustBeAny;
  }
}

export = knex;

//// [index.ts]
import "knex";
declare module "knex" {
  namespace Knex {
    function newFunc(): Knex.Interface;
  }
}




//// [index.js]
import "knex";


//// [index.d.ts]
import "knex";
module "knex" {
    namespace Knex {
        function newFunc(): Knex.Interface;
    }
}


//// [DtsFileErrors]


index.d.ts(2,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== node_modules/knex/index.d.ts (0 errors) ====
    // A bunch of random text to move the positions forward
    // A bunch of random text to move the positions forward
    // A bunch of random text to move the positions forward
    // A bunch of random text to move the positions forward
    // A bunch of random text to move the positions forward
    // A bunch of random text to move the positions forward
    // A bunch of random text to move the positions forward
    // A bunch of random text to move the positions forward
    // A bunch of random text to move the positions forward
    
    type ShouldJustBeAny = [any][0];
    
    declare namespace knex {
      export { Knex };
    }
    
    declare namespace Knex {
      interface Interface {
        method(): ShouldJustBeAny;
      }
    }
    
    export = knex;
    
==== index.d.ts (1 errors) ====
    import "knex";
    module "knex" {
    ~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        namespace Knex {
            function newFunc(): Knex.Interface;
        }
    }
    