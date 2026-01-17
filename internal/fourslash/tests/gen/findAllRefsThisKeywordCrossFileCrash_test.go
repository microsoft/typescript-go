package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

// Test for issue: Crash on find-all-references on `this` across files
// This reproduces the panic when searching for 'this' references where the container
// from one file is used to search in another file.
// Without the fix, this test panics with "slice bounds out of range".
func TestFindAllRefsThisKeywordCrossFileCrash(t *testing.T) {
	fourslash.SkipIfFailing(t)
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	
	// Create a scenario where 'this' in a class triggers cross-file search.
	// The key is having a very large file where container positions are large,
	// and a very small file that would cause slice bounds errors.
	const content = `// @Filename: largeFile.ts
namespace LargeNamespace {
    export class LargeClass {
        // Adding many properties to push container positions far into the file
        private property001: string = "long string value to increase file size and push positions forward significantly";
        private property002: string = "long string value to increase file size and push positions forward significantly";
        private property003: string = "long string value to increase file size and push positions forward significantly";
        private property004: string = "long string value to increase file size and push positions forward significantly";
        private property005: string = "long string value to increase file size and push positions forward significantly";
        private property006: string = "long string value to increase file size and push positions forward significantly";
        private property007: string = "long string value to increase file size and push positions forward significantly";
        private property008: string = "long string value to increase file size and push positions forward significantly";
        private property009: string = "long string value to increase file size and push positions forward significantly";
        private property010: string = "long string value to increase file size and push positions forward significantly";
        
        constructor() {
            /*1*/this.property001 = "init";
        }
        
        someMethod() {
            return /*2*/this;
        }
    }
}
// @Filename: a.ts
class A { m() { return /*3*/this; } }`
	
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	
	// Without the fix, this panics when finding references from marker 1 or 2
	// because it searches across files using a container with positions from largeFile
	f.VerifyBaselineFindAllReferences(t, "1", "2", "3")
}
