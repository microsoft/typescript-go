package ast

import (
	"testing"
)

func TestElementAccessExpressionText(t *testing.T) {
	factory := NewNodeFactory(NodeFactoryHooks{})
	
	// Test with a string literal argument
	stringLiteral := factory.NewStringLiteral("key")
	expression := factory.NewIdentifier("obj")
	elementAccess := factory.NewElementAccessExpression(expression, nil, stringLiteral, 0)
	
	text := elementAccess.Text()
	if text != "key" {
		t.Errorf("Expected Text() to return 'key', got '%s'", text)
	}
	
	// Test with a numeric literal argument
	numericLiteral := factory.NewNumericLiteral("123")
	elementAccess = factory.NewElementAccessExpression(expression, nil, numericLiteral, 0)
	
	text = elementAccess.Text()
	if text != "123" {
		t.Errorf("Expected Text() to return '123', got '%s'", text)
	}
	
	// Test with a non-literal argument
	nonLiteralArg := factory.NewIdentifier("nonLiteralKey")
	elementAccess = factory.NewElementAccessExpression(expression, nil, nonLiteralArg, 0)
	
	text = elementAccess.Text()
	if text != "" {
		t.Errorf("Expected Text() to return '', got '%s'", text)
	}
}