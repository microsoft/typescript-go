// @strict: true
// @target: esnext
// @noEmit: true

class Cls1 {
  get x() {
    return 1;
  }
  get x() {
    return 2;
  }
}

const obj1 = {
  get x() {
    return 1;
  },
  get x() {
    return 2;
  },
};

class Cls2 {
  get x() {
    return 1;
  }
  x = 1;
}

interface Interface1 {
  x: 1;
}
interface Interface1 {
  x: 1;
}

interface Interface2 {
  set x(x);
  get x(): number;
}
interface Interface2 {
  set x(x);
  get x(): number;
}

class Cls3 {
  set x(x) {}
  get x() {
    return 1;
  }
}
interface Cls3 {
  x: number;
}

class Cls4 {
  foo() {
    return "";
  }
  foo() {
    return "";
  }
  foo() {
    return "";
  }
}

class Cls5 {
  foo = 1;
  foo = 1;
  foo = 1;
}

class Cls6 {
  foo() {}
  foo = 1;
}

class Cls7 {
  foo = 1;
  foo() {}
}

class Cls8 {
  k = 1;
  k() {}
  k() {}
}

class Cls9 {
  k() {}
  k() {}
  k = 1;
}

class Cls10 {
  k() {}
  k = 1;
  k() {}
}

class Cls11 {
  foo = 1;

  foo(): void;
  foo(): void;
  foo() {}
}

class Cls12 {
  foo(): void;
  foo(): void;
  foo() {}

  foo = 1;
}

class Cls13 {
  get x() {
    return 1;
  }

  accessor x = 1;
}

class Cls14 {
  get x() {
    return 1;
  }

  accessor x = 1;
}

class Cls15 {
  accessor x = 1;

  get x() {
    return 1;
  }
}

class Cls16 {
  accessor x = 1;
  accessor x = 1;
}

class Cls17 {
  get x() {
    return 1;
  }
  set x(v) {}
  accessor x = 1;
}

class Cls18 {
  accessor x = 1;
  get x() {
    return 1;
  }
  set x(v) {}
}

class Cls19 {
  get x() {
    return 1;
  }
  accessor x = 1;
  set x(v) {}
}

interface Interface3 {
  set x(x);
  get x(): number;
  get x(): number;
}

interface Interface4 {
  x: number;
  get x(): number;
}

interface Interface5 {
  x: number;
  set x(v);
}

interface Interface6 {
  get x(): number;
  x: number;
}

interface Interface7 {
  set x(v);
  x: number;
}

interface Interface8 {
  x: number;
  get x(): number;
  set x(v);
}

interface Interface9 {
  get x(): number;
  set x(v);
  x: number;
}

interface Interface10 {
  get x(): number;
}
interface Interface10 {
  get x(): number;
}

interface Interface11 {
  get x(): number;
}
interface Interface11 {
  x: number;
}

interface Interface12 {
  set x(v);
}
interface Interface12 {
  x: number;
}

interface Interface13 {
  x: number;
}
interface Interface13 {
  get x(): number;
}

interface Interface14 {
  x: number;
}
interface Interface14 {
  set x(v);
}
