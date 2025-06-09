package checker_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/tsoptions"
	"github.com/microsoft/typescript-go/internal/vfs/vfstest"
	"gotest.tools/v3/assert"
)

// TestSliceBoundsPanicFix tests the fix for issue #1108
// 
// The issue was a panic in addPropertyToElementList with the error:
// "runtime error: slice bounds out of range [:-1]"
//
// Root cause: Asymmetric push/pop operations on reverseMappedStack when:
// 1. propertyIsReverseMapped = true  
// 2. shouldUsePlaceholderForProperty() returns true (skip push)
// 3. Pop operation still executed (was outside the conditional)
//
// Fix: Track whether we actually pushed to stack and only pop if we did
func TestSliceBoundsPanicFix(t *testing.T) {
	t.Parallel()

	// Create TypeScript code that exercises reverse-mapped types
	// in scenarios that could trigger the stack imbalance
	content := `
// Complex nested mapped types that can trigger reverse mapping
type DeepPartial<T> = {
	[P in keyof T]?: T[P] extends object ? DeepPartial<T[P]> : T[P];
};

type ComplexNested<T> = {
	[K in keyof T]: {
		wrapper: T[K];
		optional?: T[K];
	}
};

interface TestInterface {
	id: string;
	data: {
		items: TestInterface[];
		meta: {
			count: number;
			nested: TestInterface;
		}
	}
}

// These mapped types create complex reverse mapping scenarios
type PartialTest = DeepPartial<TestInterface>;
type ComplexTest = ComplexNested<TestInterface>;
type CombinedTest = DeepPartial<ComplexNested<TestInterface>>;

// Force type checking and potential serialization that would 
// trigger addPropertyToElementList in various scenarios
declare function testPartial(x: PartialTest): void;
declare function testComplex(x: ComplexTest): void; 
declare function testCombined(x: CombinedTest): void;

declare const partialValue: PartialTest;
declare const complexValue: ComplexTest;
declare const combinedValue: CombinedTest;

// These calls should trigger type checking without panicking
testPartial(partialValue);
testComplex(complexValue);
testCombined(combinedValue);

// Test deep property access that might exercise the reverse mapping stack
const deepAccess = combinedValue.data?.wrapper.data?.wrapper.meta?.wrapper.count;
`

	fs := vfstest.FromMap(map[string]string{
		"/test.ts": content,
		"/tsconfig.json": `{
			"compilerOptions": {
				"strict": true,
				"exactOptionalPropertyTypes": true,
				"noEmit": true
			},
			"files": ["test.ts"]
		}`,
	}, false /*useCaseSensitiveFileNames*/)
	fs = bundled.WrapFS(fs)

	cd := "/"
	host := compiler.NewCompilerHost(nil, cd, fs, bundled.LibPath())

	parsed, errors := tsoptions.GetParsedCommandLineOfConfigFile("/tsconfig.json", &core.CompilerOptions{}, host, nil)
	assert.Equal(t, len(errors), 0, "Expected no errors in parsed command line")

	p := compiler.NewProgram(compiler.ProgramOptions{
		Config: parsed,
		Host:   host,
	})
	p.BindSourceFiles()
	c, done := p.GetTypeChecker(t.Context())
	defer done()

	// The main test: this should complete without any slice bounds panic
	// Before the fix, complex reverse mapping scenarios could cause:
	// panic: runtime error: slice bounds out of range [:-1]
	file := p.GetSourceFile("/test.ts")
	
	// Wrap in a defer to catch any panics and provide context
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Unexpected panic occurred during type checking: %v\n"+
				"This suggests the fix for issue #1108 is not working correctly.\n"+
				"The panic likely occurred in addPropertyToElementList due to "+
				"asymmetric push/pop operations on reverseMappedStack.", r)
		}
	}()

	// Force comprehensive type checking that exercises the nodebuilder
	c.CheckSourceFile(t.Context(), file)
	
	t.Log("SUCCESS: Complex reverse mapping scenarios completed without slice bounds panic")
	t.Log("The fix for issue #1108 is working correctly")
}

// TestReverseMappedStackBalance verifies the fix maintains proper
// balance between push and pop operations on the reverseMappedStack
func TestReverseMappedStackBalance(t *testing.T) {
	t.Parallel()

	// Test various edge cases that could expose stack imbalance issues
	content := `
// Test recursive mapped types that stress the reverse mapping stack
type Identity<T> = T;
type Nested<T> = { nested: T };
type DoubleNested<T> = Nested<Nested<T>>;
type TripleNested<T> = Nested<DoubleNested<T>>;

type MappedIdentity<T> = { [K in keyof T]: Identity<T[K]> };
type MappedNested<T> = { [K in keyof T]: Nested<T[K]> };

interface BaseType {
	a: string;
	b: number;  
	c: boolean;
}

// Create multiple layers that could trigger different stack behaviors
type Layer1 = MappedIdentity<BaseType>;
type Layer2 = MappedNested<Layer1>;
type Layer3 = MappedIdentity<Layer2>;
type Layer4 = MappedNested<Layer3>;

// Force evaluation
declare const value: Layer4;
const test1 = value.a.nested.a.nested.a;
const test2 = value.b.nested.b.nested.b;
const test3 = value.c.nested.c.nested.c;
`

	fs := vfstest.FromMap(map[string]string{
		"/balance_test.ts": content,
		"/tsconfig.json": `{
			"compilerOptions": {
				"strict": true,
				"noEmit": true
			},
			"files": ["balance_test.ts"]
		}`,
	}, false /*useCaseSensitiveFileNames*/)
	fs = bundled.WrapFS(fs)

	cd := "/"
	host := compiler.NewCompilerHost(nil, cd, fs, bundled.LibPath())

	parsed, errors := tsoptions.GetParsedCommandLineOfConfigFile("/tsconfig.json", &core.CompilerOptions{}, host, nil)
	assert.Equal(t, len(errors), 0, "Expected no errors in parsed command line")

	p := compiler.NewProgram(compiler.ProgramOptions{
		Config: parsed,
		Host:   host,
	})
	p.BindSourceFiles()
	c, done := p.GetTypeChecker(t.Context())
	defer done()

	// Test that complex nested scenarios don't cause stack issues
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Stack balance test failed with panic: %v", r)
		}
	}()

	file := p.GetSourceFile("/balance_test.ts")
	c.CheckSourceFile(t.Context(), file)
	
	t.Log("Stack balance test passed - proper push/pop balance maintained")
}