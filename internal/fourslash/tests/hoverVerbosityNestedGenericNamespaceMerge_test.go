package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestHoverVerbosityNestedGenericNamespaceMerge(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `
declare namespace NS {
    export class Box<TElement = HTMLElement> {}
    export namespace Box {
        export interface ResultBase<TR, TJ, TN> {
            next<ARD, AJD, AND>(
                doneFilter: (
                    value: TR
                ) => ResultBase<ARD, AJD, AND> | ARD,
            ): ResultBase<ARD, AJD, AND>;
        }
        export interface Result<TR, TJ = any, TN = any> extends ResultBase<TR, TJ, TN> {}
    }
}

interface Operation<T> extends NS.Box.Result<T> {}

declare namespace Work {
    interface Item {
        run(force?: boolean): Operation<any>;
    }
}

declare const item: Work.Item;
item./*1*/run();
`

	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyBaselineHoverWithVerbosity(t, map[string][]int{"1": {0, 1, 2}})
}
