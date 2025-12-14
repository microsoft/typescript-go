// @target: esnext
// @experimentaldecorators: true

declare function dec(target: any, propertyKey: string, descriptor: PropertyDescriptor): PropertyDescriptor;

// Test case for single getter without setter
class C1 {
    @dec get accessor() { return 1; }
}

// Test case for single setter without getter
class C2 {
    @dec set accessor(value: number) { }
}

// Test case for decorated class without modifiers
declare function classDec(target: any): any;

@classDec
class C3 {
}
