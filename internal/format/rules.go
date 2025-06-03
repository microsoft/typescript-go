package format

import "github.com/microsoft/typescript-go/internal/ast"

type tokenRange struct {
	tokens     []ast.Kind
	isSpecific bool
}

type ruleSpec struct {
	leftTokenRange  tokenRange
	rightTokenRange tokenRange
	rule            Rule
}

type rule struct{}

func getAllRules() []ruleSpec {
	return nil // !!!
}
