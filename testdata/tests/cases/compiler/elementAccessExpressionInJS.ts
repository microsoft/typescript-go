// @ts-check
// @checkJs: true

// Create an object with properties
const obj = { staticKey: "value", 0: "numeric" };

// Use ElementAccessExpression with string literal (should work)
const value1 = obj["staticKey"];

// Use ElementAccessExpression with numeric literal (should work)
const value2 = obj[0];

// Use ElementAccessExpression with variable (would trigger the panic before fix)
const dynamicKey = "dynamicKey";
const value3 = obj[dynamicKey]; 

// Nested ElementAccessExpression (common in webpack code)
const nestedObj = { inner: { deep: "value" } };
const innerKey = "inner";
const deepKey = "deep";
const nestedValue = nestedObj[innerKey][deepKey];

// Create an array and access with ElementAccessExpression
const arr = ["first", "second", "third"];
const idx = 1;
const arrValue = arr[idx];