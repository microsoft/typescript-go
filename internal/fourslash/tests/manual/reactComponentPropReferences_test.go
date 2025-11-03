package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestReactComponentPropReferences(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `//@Filename: Counter.tsx
// @jsx: preserve
// @noLib: true
declare module JSX {
    interface Element { }
    interface IntrinsicElements {
    }
    interface ElementAttributesProperty { props; }
}

interface CounterProps {
    /*1*/value: number;
}

export function Counter({ value }: CounterProps) {
    return <div>{value}</div>;
}
//@Filename: App.tsx
// @jsx: preserve
// @noLib: true
import { Counter } from './Counter';

const App = () => {
    const count = 0;
    return <Counter value={count} />;
};
`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyBaselineFindAllReferences(t, "1")
}
