package checker_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/tsoptions"
	"github.com/microsoft/typescript-go/internal/vfs/vfstest"
)

// TestReverseMappedStackPanicFix tests the fix for the panic in addPropertyToElementList
// when a reverse-mapped property should use a placeholder, ensuring we don't try to pop
// from an empty reverseMappedStack.
func TestReverseMappedStackPanicFix(t *testing.T) {
	t.Parallel()

	// This test creates a scenario similar to what caused the original panic:
	// - Complex recursive mapped types
	// - Type errors that trigger type-to-string conversion 
	// - Deep nesting that could cause placeholder usage
	content := `
		// Recursive mapped type that can cause reverse mapping
		type DeepPartial<T> = {
			[P in keyof T]?: T[P] extends object ? DeepPartial<T[P]> : T[P];
		};
		
		// Complex nested interface
		interface ComplexInterface {
			id: string;
			data: {
				items: ComplexInterface[];
				meta: {
					count: number;
					nested: ComplexInterface;
				}
			}
		}
		
		// This creates deeply nested reverse mapped types
		type PartialComplex = DeepPartial<ComplexInterface>;
		
		// Force type checking errors that require type stringification
		// This is key because the original panic happened during error reporting
		declare function expectNumber(x: number): void;
		declare const value: PartialComplex;
		
		// These should trigger type errors and reverse mapping during error reporting
		expectNumber(value.id); // Error: string | undefined is not assignable to number
		expectNumber(value.data.items[0].data.meta.count); // More complex nested access
		expectNumber(value.data.meta.nested.id); // Deep nesting that could trigger placeholders
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
	if len(errors) > 0 {
		t.Fatalf("Expected no errors in parsed command line, got %d errors", len(errors))
	}

	p := compiler.NewProgram(compiler.ProgramOptions{
		Config: parsed,
		Host:   host,
	})
	p.BindSourceFiles()

	// This should NOT panic after our fix, even though it would have before
	p.CheckSourceFiles(t.Context())
	
	// If we reach here without panicking, the fix is working
	t.Log("Test completed successfully - no panic occurred")
}

// TestReverseMappedStackHandling tests edge cases around reverse mapping stack management
func TestReverseMappedStackHandling(t *testing.T) {
	t.Parallel()

	content := `
		// Create multiple levels of mapped types to stress test the stack handling
		type Level1<T> = { [K in keyof T]: T[K] };
		type Level2<T> = Level1<{ [K in keyof T]: T[K] }>;
		type Level3<T> = Level2<{ [K in keyof T]: T[K] }>;
		
		interface TestInterface {
			a: {
				b: {
					c: {
						d: string;
					}
				}
			}
		}
		
		// Deeply nested mapped type
		type DeeplyMapped = Level3<TestInterface>;
		
		// Create error scenarios that force type inspection
		declare function requireString(x: string): void;
		declare const mapped: DeeplyMapped;
		
		// Force multiple type errors that require stringification
		requireString(mapped.a.b.c.d); // Should work
		requireString(mapped.a.b.c); // Should error - object is not string
		requireString(mapped.a.b); // Should error - object is not string  
		requireString(mapped.a); // Should error - object is not string
		requireString(mapped); // Should error - object is not string
	`

	fs := vfstest.FromMap(map[string]string{
		"/edge_case.ts": content,
		"/tsconfig.json": `{
			"compilerOptions": {
				"strict": true,
				"noEmit": true
			},
			"files": ["edge_case.ts"]
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

	// This should handle the stack correctly without panicking
	p.CheckSourceFiles(t.Context())
}