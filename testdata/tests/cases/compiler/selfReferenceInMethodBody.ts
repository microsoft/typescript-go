// @noEmit: true
// @noTypesAndSymbols: true
// @strict: true

class Box<T> {
    value!: T;

    method(source: Box<Box<Box<number>>>) {
        const target: Box<Box<Box<string>>> = source;
    }
}
