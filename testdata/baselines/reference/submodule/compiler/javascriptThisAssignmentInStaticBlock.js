//// [tests/cases/compiler/javascriptThisAssignmentInStaticBlock.ts] ////

//// [a.js]
class Thing {
    static {
        this.doSomething = () => {};
    }
}

Thing.doSomething();

// GH#46468
class ElementsArray extends Array {
    static {
        const superisArray = super.isArray;
        const customIsArray = (arg)=> superisArray(arg);
        this.isArray = customIsArray;
    }
}

ElementsArray.isArray(new ElementsArray());

//// [a.js]
class Thing {
    static {
        this.doSomething = () => { };
    }
}
Thing.doSomething();
class ElementsArray extends Array {
    static {
        const superisArray = super.isArray;
        const customIsArray = (arg) => superisArray(arg);
        this.isArray = customIsArray;
    }
}
ElementsArray.isArray(new ElementsArray());
