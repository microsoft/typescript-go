package checker_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/tsoptions"
	"github.com/microsoft/typescript-go/internal/vfs/vfstest"
)

func TestReverseMappedPropertySliceBoundsPanic(t *testing.T) {
	// This test attempts to reproduce the slice bounds panic in addPropertyToElementList.
	//
	// THE ISSUE:
	// In addPropertyToElementList, there's an asymmetric push/pop operation on reverseMappedStack:
	//
	// 1. Push operation (line 2154-2156):
	//    if propertyIsReverseMapped {
	//        b.ctx.reverseMappedStack = append(b.ctx.reverseMappedStack, propertySymbol)
	//    }
	//    This only happens if !b.shouldUsePlaceholderForProperty(propertySymbol) (in else branch)
	//
	// 2. Pop operation (line 2162-2164):
	//    if propertyIsReverseMapped {
	//        b.ctx.reverseMappedStack = b.ctx.reverseMappedStack[:len(b.ctx.reverseMappedStack)-1]
	//    }
	//    This happens unconditionally if propertyIsReverseMapped is true
	//
	// PROBLEMATIC SCENARIO:
	// - propertySymbol.CheckFlags&ast.CheckFlagsReverseMapped != 0 (propertyIsReverseMapped = true)
	// - shouldUsePlaceholderForProperty(propertySymbol) returns true (skip push)
	// - Pop operation still executes because propertyIsReverseMapped is true
	// - Panic: slice bounds out of range [:-1] when stack is empty
	//
	// This test creates TypeScript code that attempts to trigger this scenario by:
	// 1. Creating reverse-mapped properties through type inference
	// 2. Creating conditions that would make shouldUsePlaceholderForProperty return true
	// 3. Forcing the nodebuilder to process these types
	
	content := `
// Attempt to create reverse-mapped properties through complex type inference
type MappedType<T> = { [K in keyof T]: T[K] };

// Conditional type that forces reverse mapping during inference
type ExtractMapped<T> = T extends MappedType<infer U> ? U : never;

// Create nested interface to trigger deep processing
interface TargetInterface {
  prop1: string;
  prop2: {
    nested1: number;
    nested2: {
      deep: boolean;
    }
  };
}

// Recursive mapped type to populate reverse mapping stack
type RecursiveMapped<T> = {
  [K in keyof T]: T[K] extends object ? RecursiveMapped<T[K]> : T[K]
};

// Create scenario with reverse mapping inference
type InferredType = ExtractMapped<MappedType<TargetInterface>>;

// Create deep nesting to trigger placeholder logic
type DeepType = RecursiveMapped<TargetInterface>;

// Combine types to create complex reverse mapping scenario
type ComplexType = InferredType & DeepType;

// Force type processing
declare const testVariable: ComplexType;

// Additional patterns that might trigger reverse mapping
type ConditionalMapped<T> = T extends object ? {
  [K in keyof T]: T[K] extends MappedType<infer R> ? R : T[K]
} : never;

interface ComplexInterface {
  data: MappedType<TargetInterface>;
  recursive: RecursiveMapped<TargetInterface>;
}

type ProcessedComplex = ConditionalMapped<ComplexInterface>;
declare const complexTest: ProcessedComplex;
`

	fs := vfstest.FromMap(map[string]string{
		"/test.ts": content,
		"/tsconfig.json": `{
			"compilerOptions": {
				"strict": true,
				"target": "es2015"
			},
			"files": ["test.ts"]
		}`,
	}, false)
	fs = bundled.WrapFS(fs)

	cd := "/"
	host := compiler.NewCompilerHost(nil, cd, fs, bundled.LibPath())

	parsed, errors := tsoptions.GetParsedCommandLineOfConfigFile("/tsconfig.json", &core.CompilerOptions{}, host, nil)
	if len(errors) > 0 {
		t.Fatalf("Expected no errors in parsed command line, got: %v", errors)
	}

	p := compiler.NewProgram(compiler.ProgramOptions{
		Config: parsed,
		Host:   host,
	})
	
	p.BindSourceFiles()
	
	// Set up panic recovery
	defer func() {
		if r := recover(); r != nil {
			// Check for the specific slice bounds panic
			switch v := r.(type) {
			case string:
				if v == "runtime error: slice bounds out of range [:-1]" {
					t.Logf("SUCCESS: Reproduced the slice bounds panic: %v", r)
					return
				}
			case interface{ Error() string }:
				errStr := v.Error()
				if errStr == "runtime error: slice bounds out of range [:-1]" {
					t.Logf("SUCCESS: Reproduced the slice bounds panic: %v", r)
					return
				}
			}
			// Different panic - re-throw
			panic(r)
		}
		
		// Test didn't reproduce the panic
		t.Fatal("Expected slice bounds panic [:-1] but test completed without panic")
	}()
	
	// Process the TypeScript code to attempt triggering the nodebuilder
	c, done := p.GetTypeChecker(t.Context())
	defer done()
	
	file := p.GetSourceFile("/test.ts")
	if file == nil {
		t.Fatal("Could not get source file")
	}
	
	// Force type checking and nodebuilder usage
	for _, stmt := range file.Statements.Nodes {
		switch stmt.Kind {
		case 260: // TypeAliasDeclaration
			if symbol := c.GetSymbolAtLocation(stmt.Name()); symbol != nil {
				if typeOfSymbol := c.GetTypeOfSymbolAtLocation(symbol, stmt); typeOfSymbol != nil {
					// Force nodebuilder processing through TypeToString
					_ = c.TypeToString(typeOfSymbol)
				}
			}
		case 240: // VariableStatement  
			varStmt := stmt.AsVariableStatement()
			for _, decl := range varStmt.DeclarationList.AsVariableDeclarationList().Declarations.Nodes {
				if symbol := c.GetSymbolAtLocation(decl.Name()); symbol != nil {
					if typeOfSymbol := c.GetTypeOfSymbolAtLocation(symbol, decl); typeOfSymbol != nil {
						// Force nodebuilder processing through TypeToString
						_ = c.TypeToString(typeOfSymbol)
					}
				}
			}
		}
	}
}