//// [tests/cases/compiler/super1.ts] ////

//// [super1.ts]
// Case 1
class Base1 {
    public foo() {
        return "base";
    }
}

class Sub1 extends Base1 {
    public bar() {
        return "base";
    }
}

class SubSub1 extends Sub1 {
    public bar() {
        return super.super.foo;
    }
}

// Case 2
class Base2 {
    public foo() {
        return "base";
    }
}

class SubE2 extends Base2 {
    public bar() {
        return super.prototype.foo = null;
    }
}

// Case 3
class Base3 {
    public foo() {
        return "base";
    }
}

class SubE3 extends Base3 {
    public bar() {
        return super.bar();
    }
}

// Case 4
module Base4 {
    class Sub4 {
        public x(){
            return "hello";
        }
    }
    
    export class SubSub4 extends Sub4{
        public x(){
            return super.x();
        }
    }
    
    export class Sub4E {
        public x() {
            return super.x();
        }
    }
}


//// [super1.js]
// Case 1
class Base1 {
    foo() {
        return "base";
    }
}
class Sub1 extends Base1 {
    bar() {
        return "base";
    }
}
class SubSub1 extends Sub1 {
    bar() {
        return super.super.foo;
    }
}
// Case 2
class Base2 {
    foo() {
        return "base";
    }
}
class SubE2 extends Base2 {
    bar() {
        return super.prototype.foo = null;
    }
}
// Case 3
class Base3 {
    foo() {
        return "base";
    }
}
class SubE3 extends Base3 {
    bar() {
        return super.bar();
    }
}
// Case 4
var Base4;
(function (Base4) {
    class Sub4 {
        x() {
            return "hello";
        }
    }
    class SubSub4 extends Sub4 {
        x() {
            return super.x();
        }
    }
    Base4.SubSub4 = SubSub4;
    class Sub4E {
        x() {
            return super.x();
        }
    }
    Base4.Sub4E = Sub4E;
})(Base4 || (Base4 = {}));
