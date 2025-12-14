// @target: esnext
// @experimentaldecorators: true

// Test case for decorated class without modifiers
declare function classDec(target: any): any;

@classDec
class C {
}
