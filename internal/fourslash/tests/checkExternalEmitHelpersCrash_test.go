package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestCheckExternalEmitHelpersCrash(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `
// @Filename: /tsconfig.json
{
    "compilerOptions": {
        "target": "es2015",
        "module": "commonjs",
        "importHelpers": true
    }
}

// @Filename: /node_modules/tslib/package.json
{ "name": "tslib", "main": "tslib.js", "typings": "tslib.d.ts" }

// @Filename: /node_modules/tslib/tslib.d.ts
export declare function __awaiter(thisArg: any, _arguments: any, P: Function, generator: Function): any;
export declare function __decorate(decorators: Function[], target: any, key?: string | symbol, desc?: any): any;
export declare function __esDecorate(ctor: any, descriptorIn: any, decorators: any, contextIn: any, initializers: any, extraInitializers: any): void;
export declare function __runInitializers(thisArg: any, initializers: any, value?: any): any;
export declare function __setFunctionName(f: any, name: any, prefix?: string): any;

// @Filename: /main.ts
export async function doStuff() {
    return 1;
}

function decorator(target: any, context: any) {}

export
@decorator
class /*1*/MyClass {
}
`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyQuickInfoAt(t, "1", "class MyClass", "")
}
