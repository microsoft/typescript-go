//// [tests/cases/conformance/es6/computedProperties/computedPropertyNames52.ts] ////

//// [computedPropertyNames52.js]
const array = [];
for (let i = 0; i < 10; ++i) {
    array.push(class C {
        [i] = () => C;
        static [i] = 100;
    })
}


//// [computedPropertyNames52.js]
const array = [];
for (let i = 0; i < 10; ++i) {
    array.push(class C {
        [i] = () => C;
        static [i] = 100;
    });
}
