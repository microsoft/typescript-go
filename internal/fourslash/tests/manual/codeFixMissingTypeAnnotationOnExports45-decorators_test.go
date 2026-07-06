package slash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/slash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestCodeFixMissingTypeAnnotationOnExports45_decorators(t *testing.T) {
	t.Parallel()
	      testutil.RecoverAndFail(t, "Panic on fourslash test")
	      content = `// @isolatedDeclarations: true
// @declaration: false
// @Filename: /code.ts
function classDecorator<T extends Function> (value: T, context: ClassDecoratorContext) {}
function methodDecorator<This> (
  target: (...args: number[])=> number,
  context: ClassMethodDecoratorContext<This, (this: This, ...args: number[]) => number>) {}
function getterDecorator(value: Function, context: ClassGetterDecoratorContext) {}
function setterDecorator(value: Function, context: ClassSetterDecoratorContext) {}
function fieldDecorator(value: undefined, context: ClassFieldDecoratorContext) {}
function foobar() { return 42;}

@classDecorator
export class A {
  @methodDecorator
  sum(...args: number[]) {
    return args.reduce((a, b) => a + b, 0);
  }
  getSelf() {
    return this;
  }
  @getterDecorator
  get a() {
    return foobar();
  }
  @setterDecorator
  set a(value) {}

  @fieldDecorator classProp = foobar();
}`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	     done()
	f.VerifyCodeFixAll(t, fourslash.VerifyCodeFixAllOptions{
		FixID: "fixMissingTypeAnnotationOnExports",
		NewFileContent: `function classDecorator<T extends Function> (value: T, context: ClassDecoratorContext) {}
function methodDecorator<This> (
  target: (...args: number[])=> number,
  context: ClassMethodDecoratorContext<This, (this: This, ...args: number[]) => number>) {}
function getterDecorator(value: Function, context: ClassGetterDecoratorContext) {}
function setterDecorator(value: Function, context: ClassSetterDecoratorContext) {}
function fieldDecorator(value: undefined, context: ClassFieldDecoratorContext) {}
function foobar() { return 42;}

@classDecorator
export class A {
  @methodDecorator
  sum(...args: number[]): number {
    return args.reduce((a, b) => a + b, 0);
  }
  getSelf(): this {
    return this;
  }
  @getterDecorator
  get a(): number {
    return foobar();
  }
  @setterDecorator
  set a(value: number) {}

  @fieldDecorator classProp: number = foobar();
}`,
	})
}
