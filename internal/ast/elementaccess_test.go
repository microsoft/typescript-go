package ast

import (
	"testing"
)

func TestGetElementOrPropertyAccessName(t *testing.T) {
	factory := NewNodeFactory(NodeFactoryHooks{})
	
	// Test with a string literal argument
	stringLiteral := factory.NewStringLiteral("key")
	expression := factory.NewIdentifier("obj")
	elementAccess := factory.NewElementAccessExpression(expression, nil, stringLiteral, 0)
	
	name := GetElementOrPropertyAccessName(elementAccess)
	if name != "key" {
		t.Errorf("Expected GetElementOrPropertyAccessName to return 'key', got '%s'", name)
	}
	
	// Test with a numeric literal argument
	numericLiteral := factory.NewNumericLiteral("123")
	elementAccess = factory.NewElementAccessExpression(expression, nil, numericLiteral, 0)
	
	name = GetElementOrPropertyAccessName(elementAccess)
	if name != "123" {
		t.Errorf("Expected GetElementOrPropertyAccessName to return '123', got '%s'", name)
	}
	
	// Test with a non-literal argument
	nonLiteralArg := factory.NewIdentifier("nonLiteralKey")
	elementAccess = factory.NewElementAccessExpression(expression, nil, nonLiteralArg, 0)
	
	name = GetElementOrPropertyAccessName(elementAccess)
	if name != "" {
		t.Errorf("Expected GetElementOrPropertyAccessName to return '', got '%s'", name)
	}
}