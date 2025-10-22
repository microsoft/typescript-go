package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestDeclarationMapGoToDefinitionChanges(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @Filename: index.ts
export class Foo {
    member: string;
    methodName(propName: SomeType): void {}
    otherMethod() {
        if (Math.random() > 0.5) {
            return {x: 42};
        }
        return {y: "yes"};
    }
}

export interface SomeType {
    member: number;
}
// @Filename: index.d.ts.map
{"version":3,"file":"index.d.ts","sourceRoot":"","sources":["index.ts"],"names":[],"mappings":"AAAA,qBAAa,GAAG;IACZ,MAAM,EAAE,MAAM,CAAC;IACV,UAAU,CAAC,QAAQ,EAAE,QAAQ,GAAG,IAAI;CAC5C;AAED,MAAM,WAAW,QAAQ;IACrB,MAAM,EAAE,MAAM,CAAC;CAClB"}
// @Filename: index.d.ts
export declare class Foo {
    member: string;
    methodName(propName: SomeType): void;
}
export interface SomeType {
    member: number;
}
//# sourceMappingURL=index.d.ts.map
// @Filename: mymodule.ts
import * as mod from "./index";
const instance = new mod.Foo();
instance.[|/*1*/methodName|]({member: 12});
instance.[|/*2*/otherMethod|]();`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyBaselineGoToDefinition(t, "1")
	f.GoToFile(t, "index.d.ts")
	f.ReplaceAll(t, `export declare class Foo {
    member: string;
    methodName(propName: SomeType): void;
    otherMethod(): {
        x: number;
        y?: undefined;
    } | {
        y: string;
        x?: undefined;
    };
}
export interface SomeType {
    member: number;
}
//# sourceMappingURL=index.d.ts.map`)
	f.GoToFile(t, "index.d.ts.map")
	f.ReplaceAll(t, `{"version":3,"file":"index.d.ts","sourceRoot":"","sources":["index.ts"],"names":[],"mappings":"AAAA,qBAAa,GAAG;IACZ,MAAM,EAAE,MAAM,CAAC;IACV,UAAU,CAAC,QAAQ,EAAE,QAAQ,GAAG,IAAI;IACzC,WAAW;;;;;;;CAMd;AAED,MAAM,WAAW,QAAQ;IACrB,MAAM,EAAE,MAAM,CAAC;CAClB"}`)
	f.VerifyBaselineGoToDefinition(t, "2")
}
