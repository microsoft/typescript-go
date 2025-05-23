//// [tests/cases/compiler/intersectionSatisfiesConstraint.ts] ////

//// [intersectionSatisfiesConstraint.ts]
interface FirstInterface {
    commonProperty: number
}

interface SecondInterface {
    commonProperty: number
}

const myFirstFunction = <T extends FirstInterface | SecondInterface>(param1: T) => {
    const newParam: T & { otherProperty: number } = Object.assign(param1, { otherProperty: 3 })
    mySecondFunction(newParam)
}

const mySecondFunction = <T extends { commonProperty: number, otherProperty: number }>(newParam: T) => {
    return newParam
}


//// [intersectionSatisfiesConstraint.js]
const myFirstFunction = (param1) => {
    const newParam = Object.assign(param1, { otherProperty: 3 });
    mySecondFunction(newParam);
};
const mySecondFunction = (newParam) => {
    return newParam;
};
