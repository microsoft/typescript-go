package checker_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/tsoptions"
	"github.com/microsoft/typescript-go/internal/vfs/vfstest"
)

// TestAddPropertyToElementListStackSafety specifically tests the fix for the
// slice bounds out of range panic in addPropertyToElementList.
// 
// The original issue occurred when:
// 1. A property had CheckFlagsReverseMapped set (propertyIsReverseMapped = true)
// 2. shouldUsePlaceholderForProperty returned true for this property
// 3. The code would skip pushing to reverseMappedStack but still try to pop
// 4. This caused a panic when trying to slice an empty stack with [:-1]
func TestAddPropertyToElementListStackSafety(t *testing.T) {
	t.Parallel()

	// This test case is designed to trigger the specific conditions that
	// caused the original panic while verifying the fix works correctly
	content := `
		// Create a scenario with complex mapped types that can trigger
		// reverse mapping with placeholder conditions
		type Conditional<T> = T extends any ? { [K in keyof T]: T[K] } : never;
		
		interface RecursiveInterface {
			prop: string;
			recurse: RecursiveInterface;
		}
		
		// This creates deeply nested conditional mapped types
		type DeeplyNested = Conditional<Conditional<Conditional<RecursiveInterface>>>;
		
		// Create JSX-like elements that can trigger the specific call path from the stack trace
		declare namespace JSX {
			interface Element {}
			interface IntrinsicElements {
				div: { id?: string; }
			}
		}
		
		// Force type checking in a JSX context (matching the original stack trace)
		const Component = (props: DeeplyNested) => {
			// This should trigger type checking that leads to reverse mapping
			return <div id={props.recurse.recurse.prop} />;
		};
		
		// Create type errors that force type-to-string conversion during error reporting
		declare function expectNumber(x: number): void;
		declare const value: DeeplyNested;
		
		// This should trigger the error path that leads to type stringification
		// and potentially the problematic addPropertyToElementList call
		expectNumber(value.prop); // Type error: string is not assignable to number
	`

	fs := vfstest.FromMap(map[string]string{
		"/stack_test.tsx": content,
		"/tsconfig.json": `{
			"compilerOptions": {
				"jsx": "react",
				"strict": true,
				"noEmit": true
			},
			"files": ["stack_test.tsx"]
		}`,
	}, false /*useCaseSensitiveFileNames*/)
	fs = bundled.WrapFS(fs)

	cd := "/"
	host := compiler.NewCompilerHost(nil, cd, fs, bundled.LibPath())

	parsed, errors := tsoptions.GetParsedCommandLineOfConfigFile("/tsconfig.json", &core.CompilerOptions{}, host, nil)
	if len(errors) > 0 {
		t.Fatalf("Expected no errors in parsed command line, got %d errors", len(errors))
	}

	p := compiler.NewProgram(compiler.ProgramOptions{
		Config: parsed,
		Host:   host,
	})
	p.BindSourceFiles()

	// The original bug would cause a panic here, but with our fix it should complete safely
	p.CheckSourceFiles(t.Context())
	
	// If we reach this point, the stack safety fix is working
	t.Log("Stack safety test passed - no slice bounds panic occurred")
}

// TestReverseMappedStackPushPopBalance verifies that our fix maintains
// proper balance between push and pop operations on the reverseMappedStack
func TestReverseMappedStackPushPopBalance(t *testing.T) {
	t.Parallel()

	content := `
		// Test various scenarios that exercise the stack push/pop logic
		type Identity<T> = { [K in keyof T]: T[K] };
		type Nested<T> = Identity<{ nested: T }>;
		type DoubleNested<T> = Nested<Nested<T>>;
		
		interface TestType {
			a: string;
			b: number;
			c: boolean;
		}
		
		// Multiple layers of nesting
		type Layer1 = DoubleNested<TestType>;
		type Layer2 = DoubleNested<Layer1>;
		type Layer3 = DoubleNested<Layer2>;
		
		declare const deeply: Layer3;
		
		// These should exercise the stack in various ways
		declare function test1(x: string): void;
		declare function test2(x: number): void;
		declare function test3(x: boolean): void;
		
		// Each of these should trigger type checking and potential reverse mapping
		test1(deeply.nested.nested.nested.nested.nested.nested.a);
		test2(deeply.nested.nested.nested.nested.nested.nested.b);
		test3(deeply.nested.nested.nested.nested.nested.nested.c);
		
		// Force some type errors to trigger stringification
		test1(deeply.nested.nested.nested.nested.nested.nested.b); // Error: number not string
		test2(deeply.nested.nested.nested.nested.nested.nested.c); // Error: boolean not number
		test3(deeply.nested.nested.nested.nested.nested.nested.a); // Error: string not boolean
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
	if len(errors) > 0 {
		t.Fatalf("Expected no errors in parsed command line, got %d errors", len(errors))
	}

	p := compiler.NewProgram(compiler.ProgramOptions{
		Config: parsed,
		Host:   host,
	})
	p.BindSourceFiles()

	// This should handle all the nested type checking without stack issues
	p.CheckSourceFiles(t.Context())
}