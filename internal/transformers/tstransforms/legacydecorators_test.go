package tstransforms_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/binder"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/printer"
	"github.com/microsoft/typescript-go/internal/testutil/emittestutil"
	"github.com/microsoft/typescript-go/internal/testutil/parsetestutil"
	"github.com/microsoft/typescript-go/internal/transformers"
	"github.com/microsoft/typescript-go/internal/transformers/tstransforms"
)

// decorateHelper is the expected output for the __decorate helper
const decorateHelper = `var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
`

func TestLegacyDecoratorsTransformer(t *testing.T) {
	t.Parallel()
	data := []struct {
		title  string
		input  string
		output string
	}{
		// Test case for single getter without setter (secondAccessor is nil)
		{
			title: "single getter with decorator",
			input: `declare function dec(target: any, propertyKey: string, descriptor: PropertyDescriptor): PropertyDescriptor;
class C {
    @dec get accessor() { return 1; }
}`,
			output: decorateHelper + `class C {
    get accessor() { return 1; }
}
__decorate([
    dec
], C.prototype, "accessor", null);`,
		},
		// Test case for single setter without getter
		{
			title: "single setter with decorator",
			input: `declare function dec(target: any, propertyKey: string, descriptor: PropertyDescriptor): PropertyDescriptor;
class C {
    @dec set accessor(value: number) { }
}`,
			output: decorateHelper + `class C {
    set accessor(value) { }
}
__decorate([
    dec
], C.prototype, "accessor", null);`,
		},
		// Test case for decorated class without modifiers (empty modifier list)
		{
			title: "decorated class without modifiers",
			input: `declare function dec(target: any): any;
@dec
class C {
}`,
			output: decorateHelper + `let C = class C {
};
C = __decorate([
    dec
], C);`,
		},
	}

	for _, rec := range data {
		t.Run(rec.title, func(t *testing.T) {
			t.Parallel()
			options := &core.CompilerOptions{}
			file := parsetestutil.ParseTypeScript(rec.input, false /*jsx*/)
			parsetestutil.CheckDiagnostics(t, file)
			binder.BindSourceFile(file)
			emitContext := printer.NewEmitContext()
			resolver := binder.NewReferenceResolver(options, binder.ReferenceResolverHooks{})
			opts := &transformers.TransformOptions{Context: emitContext, CompilerOptions: options, Resolver: resolver}
			file = tstransforms.NewTypeEraserTransformer(opts).TransformSourceFile(file)
			file = tstransforms.NewLegacyDecoratorsTransformer(opts).TransformSourceFile(file)
			emittestutil.CheckEmit(t, emitContext, file, rec.output)
		})
	}
}
