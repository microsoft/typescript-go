package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestImportFixJsxNoCrash(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @jsx: react-jsx
// @module: commonjs
// @Filename: /component.tsx
export function Component(props: { children?: any }) { return <div>{props.children}</div>; }
// @Filename: /index.tsx
const App = () => {
    return (
        <Component/**/ />
    );
};`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.GoToMarker(t, "")
	f.VerifyImportFixAtPosition(t, []string{
		`import { Component } from "./component";

const App = () => {
    return (
        <Component />
    );
};`,
	}, nil /*preferences*/)
}
