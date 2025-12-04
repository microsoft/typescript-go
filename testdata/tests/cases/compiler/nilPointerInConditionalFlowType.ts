// @allowJs: true
// @checkJs: true
// @declaration: true
// @outDir: ./out
// @filename: test.js
class MyClass {
    /**
     * @param {string | undefined} param1
     * @param {number} param2
     * @param {boolean} param3
     */
    myMethod(param1, param2, param3) {
        if (param1) {
            console.log(param1);
        }
    }
}
