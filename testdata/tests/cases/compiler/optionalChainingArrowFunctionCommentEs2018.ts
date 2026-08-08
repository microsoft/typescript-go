// @target: es2018

const thing = { nested: { condition: true } };
const wat = () =>
    // explanatory comment
    thing?.nested?.condition ? "pass" : "fail";
