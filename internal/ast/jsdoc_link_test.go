// Test to verify JSDoc link processing doesn't panic

package ast

import (
	"testing"
)

func TestJSDocLinkText(t *testing.T) {
	factory := NewNodeFactory(NodeFactoryHooks{})
	
	// Create JSDoc links with text
	text := []string{"some", " text"}
	
	// Test JSDocLink
	jsDocLink := factory.NewJSDocLink(nil, text)
	result := jsDocLink.Text()
	expected := "some text"
	if result != expected {
		t.Errorf("JSDocLink.Text() = %q, want %q", result, expected)
	}
	
	// Test JSDocLinkCode
	jsDocLinkCode := factory.NewJSDocLinkCode(nil, text)
	result = jsDocLinkCode.Text()
	if result != expected {
		t.Errorf("JSDocLinkCode.Text() = %q, want %q", result, expected)
	}
	
	// Test JSDocLinkPlain
	jsDocLinkPlain := factory.NewJSDocLinkPlain(nil, text)
	result = jsDocLinkPlain.Text()
	if result != expected {
		t.Errorf("JSDocLinkPlain.Text() = %q, want %q", result, expected)
	}
}