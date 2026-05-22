package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestTsxFindAllReferences11VS(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `//@Filename: file.tsx
// @jsx: preserve
// @noLib: true
declare namespace JSX {
    interface Element { }
    interface IntrinsicElements {
    }
    interface ElementAttributesProperty { props; }
}
interface ClickableProps {
    children?: string;
    className?: string;
}
interface ButtonProps extends ClickableProps {
    onClick(event?: React.MouseEvent<HTMLButtonElement>): void;
}
interface LinkProps extends ClickableProps {
    goTo: string;
}
declare function MainButton(buttonProps: ButtonProps): JSX.Element;
declare function MainButton(linkProps: LinkProps): JSX.Element;
declare function MainButton(props: ButtonProps | LinkProps): JSX.Element;
let opt = <MainButton /*1*/wrong />;`
	f, done := fourslash.NewFourslash(t, &lsproto.ClientCapabilities{VSSupportsVisualStudioExtensions: new(true)}, content)
	defer done()
	f.VerifyBaselineVsFindAllReferences(t, "1")
}
