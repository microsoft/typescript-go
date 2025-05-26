package ls_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil/lstestutil"
)

const content = `export {};
interface Point {
    x: number;
    y: number;
}
declare const p: Point;
p./*a*/`

func TestBasicInterfaceMembers(t *testing.T) {
	cap := &lsproto.ClientCapabilities{}
	f, done := lstestutil.NewFourslash(t, cap, content, "basicInterfaceMembers.ts")
	f.VerifyCompletions(t, "a", nil)
	done()
	// !!! roughly:
	// !!! parse content into multiple files with markers and ranges etc
	// !!! create an ls
	// !!! call ls method
	// !!! verify result
}
