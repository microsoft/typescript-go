//// [tests/cases/compiler/declFileEmitDeclarationOnly.ts] ////

//// [helloworld.ts]
const Log = {
  info(msg: string) {}
}

class HelloWorld {
  constructor(private name: string) {
  }

  public hello() {
    Log.info(`Hello ${this.name}`);
  }
}




//// [helloworld.d.ts]
const Log: {
    info(msg: string): void;
};
class HelloWorld {
    private name;
    constructor(name: string);
    hello(): void;
}


//// [DtsFileErrors]


helloworld.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== helloworld.d.ts (1 errors) ====
    const Log: {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        info(msg: string): void;
    };
    class HelloWorld {
        private name;
        constructor(name: string);
        hello(): void;
    }
    