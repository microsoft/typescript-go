package checker

import (
	"reflect"
	"unsafe"

	"github.com/microsoft/typescript-go/internal/ast"
)

func (c *Checker) GetAliasVariancesForTesting(symbol *ast.Symbol) []VarianceFlags {
	return c.getAliasVariances(symbol)
}

func (c *Checker) MarkVarianceInProgressForTesting(symbol *ast.Symbol) func() {
	links := c.varianceLinks.Get(symbol)
	oldVariances := links.variances
	oldInVarianceComputation := c.inVarianceComputation

	links.variances = []VarianceFlags{}
	c.inVarianceComputation = true

	var oldVarianceStack reflect.Value
	var varianceStack reflect.Value
	if stack, ok := varianceStackForTesting(c); ok {
		varianceStack = stack
		oldVarianceStack = reflect.MakeSlice(stack.Type(), stack.Len(), stack.Len())
		reflect.Copy(oldVarianceStack, stack)
		stack.Set(reflect.Append(stack, reflect.ValueOf(symbol)))
	}

	return func() {
		links.variances = oldVariances
		c.inVarianceComputation = oldInVarianceComputation
		if varianceStack.IsValid() {
			varianceStack.Set(oldVarianceStack)
		}
	}
}

func varianceStackForTesting(c *Checker) (reflect.Value, bool) {
	field := reflect.ValueOf(c).Elem().FieldByName("varianceStack")
	if !field.IsValid() {
		return reflect.Value{}, false
	}
	return reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem(), true
}
