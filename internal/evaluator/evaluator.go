package evaluator

import (
	"strings"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/binder"
	"github.com/microsoft/typescript-go/internal/compiler/diagnostics"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/jsnum"
)

type EvaluatorResult struct {
	Value                 any
	IsSyntacticallyString bool
	ResolvedOtherFiles    bool
	HasExternalReferences bool
}

type Evaluator func(expr *ast.Node, location *ast.Node) EvaluatorResult

func newEvaluator(evaluateEntity Evaluator) Evaluator {
	var evaluate Evaluator
	evaluateTemplateExpression := func(expr *ast.Node, location *ast.Node) EvaluatorResult {
		var sb strings.Builder
		sb.WriteString(expr.AsTemplateExpression().Head.Text())
		resolvedOtherFiles := false
		hasExternalReferences := false
		for _, span := range expr.AsTemplateExpression().TemplateSpans.Nodes {
			spanResult := evaluate(span.Expression(), location)
			if spanResult.Value == nil {
				return EvaluatorResult{nil, true /*isSyntacticallyString*/, false, false}
			}
			sb.WriteString(anyToString(spanResult.Value))
			sb.WriteString(span.AsTemplateSpan().Literal.Text())
			resolvedOtherFiles = resolvedOtherFiles || spanResult.ResolvedOtherFiles
			hasExternalReferences = hasExternalReferences || spanResult.HasExternalReferences
		}
		return EvaluatorResult{sb.String(), true, resolvedOtherFiles, hasExternalReferences}
	}
	evaluate = func(expr *ast.Node, location *ast.Node) EvaluatorResult {
		isSyntacticallyString := false
		resolvedOtherFiles := false
		hasExternalReferences := false
		// It's unclear when/whether we should consider skipping other kinds of outer expressions.
		// Type assertions intentionally break evaluation when evaluating literal types, such as:
		//     type T = `one ${"two" as any} three`; // string
		// But it's less clear whether such an assertion should break enum member evaluation:
		//     enum E {
		//       A = "one" as any
		//     }
		// SatisfiesExpressions and non-null assertions seem to have even less reason to break
		// emitting enum members as literals. However, these expressions also break Babel's
		// evaluation (but not esbuild's), and the isolatedModules errors we give depend on
		// our evaluation results, so we're currently being conservative so as to issue errors
		// on code that might break Babel.
		expr = ast.SkipParentheses(expr)
		switch expr.Kind {
		case ast.KindPrefixUnaryExpression:
			result := evaluate(expr.AsPrefixUnaryExpression().Operand, location)
			resolvedOtherFiles = result.ResolvedOtherFiles
			hasExternalReferences = result.HasExternalReferences
			if value, ok := result.Value.(jsnum.Number); ok {
				switch expr.AsPrefixUnaryExpression().Operator {
				case ast.KindPlusToken:
					return EvaluatorResult{value, isSyntacticallyString, resolvedOtherFiles, hasExternalReferences}
				case ast.KindMinusToken:
					return EvaluatorResult{-value, isSyntacticallyString, resolvedOtherFiles, hasExternalReferences}
				case ast.KindTildeToken:
					return EvaluatorResult{value.BitwiseNOT(), isSyntacticallyString, resolvedOtherFiles, hasExternalReferences}
				}
			}
		case ast.KindBinaryExpression:
			left := evaluate(expr.AsBinaryExpression().Left, location)
			right := evaluate(expr.AsBinaryExpression().Right, location)
			operator := expr.AsBinaryExpression().OperatorToken.Kind
			isSyntacticallyString = (left.IsSyntacticallyString || right.IsSyntacticallyString) && expr.AsBinaryExpression().OperatorToken.Kind == ast.KindPlusToken
			resolvedOtherFiles = left.ResolvedOtherFiles || right.ResolvedOtherFiles
			hasExternalReferences = left.HasExternalReferences || right.HasExternalReferences
			leftNum, leftIsNum := left.Value.(jsnum.Number)
			rightNum, rightIsNum := right.Value.(jsnum.Number)
			if leftIsNum && rightIsNum {
				switch operator {
				case ast.KindBarToken:
					return EvaluatorResult{leftNum.BitwiseOR(rightNum), isSyntacticallyString, resolvedOtherFiles, hasExternalReferences}
				case ast.KindAmpersandToken:
					return EvaluatorResult{leftNum.BitwiseAND(rightNum), isSyntacticallyString, resolvedOtherFiles, hasExternalReferences}
				case ast.KindGreaterThanGreaterThanToken:
					return EvaluatorResult{leftNum.SignedRightShift(rightNum), isSyntacticallyString, resolvedOtherFiles, hasExternalReferences}
				case ast.KindGreaterThanGreaterThanGreaterThanToken:
					return EvaluatorResult{leftNum.UnsignedRightShift(rightNum), isSyntacticallyString, resolvedOtherFiles, hasExternalReferences}
				case ast.KindLessThanLessThanToken:
					return EvaluatorResult{leftNum.LeftShift(rightNum), isSyntacticallyString, resolvedOtherFiles, hasExternalReferences}
				case ast.KindCaretToken:
					return EvaluatorResult{leftNum.BitwiseXOR(rightNum), isSyntacticallyString, resolvedOtherFiles, hasExternalReferences}
				case ast.KindAsteriskToken:
					return EvaluatorResult{leftNum * rightNum, isSyntacticallyString, resolvedOtherFiles, hasExternalReferences}
				case ast.KindSlashToken:
					return EvaluatorResult{leftNum / rightNum, isSyntacticallyString, resolvedOtherFiles, hasExternalReferences}
				case ast.KindPlusToken:
					return EvaluatorResult{leftNum + rightNum, isSyntacticallyString, resolvedOtherFiles, hasExternalReferences}
				case ast.KindMinusToken:
					return EvaluatorResult{leftNum - rightNum, isSyntacticallyString, resolvedOtherFiles, hasExternalReferences}
				case ast.KindPercentToken:
					return EvaluatorResult{leftNum.Remainder(rightNum), isSyntacticallyString, resolvedOtherFiles, hasExternalReferences}
				case ast.KindAsteriskAsteriskToken:
					return EvaluatorResult{leftNum.Exponentiate(rightNum), isSyntacticallyString, resolvedOtherFiles, hasExternalReferences}
				}
			}
			leftStr, leftIsStr := left.Value.(string)
			rightStr, rightIsStr := right.Value.(string)
			if (leftIsStr || leftIsNum) && (rightIsStr || rightIsNum) && operator == ast.KindPlusToken {
				if leftIsNum {
					leftStr = leftNum.String()
				}
				if rightIsNum {
					rightStr = rightNum.String()
				}
				return EvaluatorResult{leftStr + rightStr, isSyntacticallyString, resolvedOtherFiles, hasExternalReferences}
			}
		case ast.KindStringLiteral, ast.KindNoSubstitutionTemplateLiteral:
			return EvaluatorResult{expr.Text(), true /*isSyntacticallyString*/, false, false}
		case ast.KindTemplateExpression:
			return evaluateTemplateExpression(expr, location)
		case ast.KindNumericLiteral:
			return EvaluatorResult{jsnum.FromString(expr.Text()), false, false, false}
		case ast.KindIdentifier, ast.KindElementAccessExpression:
			return evaluateEntity(expr, location)
		case ast.KindPropertyAccessExpression:
			if ast.IsEntityNameExpression(expr) {
				return evaluateEntity(expr, location)
			}
		}
		return EvaluatorResult{nil, isSyntacticallyString, resolvedOtherFiles, hasExternalReferences}
	}
	return evaluate
}

func anyToString(v any) string {
	// !!! This function should behave identically to the expression `"" + v` in JS
	switch v := v.(type) {
	case string:
		return v
	case jsnum.Number:
		return v.String()
	case bool:
		return core.IfElse(v, "true", "false")
	case jsnum.PseudoBigInt:
		return v.String()
	}
	panic("Unhandled case in anyToString")
}

type ConstantEvaluator struct {
	evaluate                           Evaluator
	CompilerOptions                    *core.CompilerOptions
	NameResolver                       *binder.NameResolver
	GetEnumDeclarationValuesComputed   func(node *ast.Node) bool                                                               // Required
	SetEnumDeclarationValuesComputed   func(node *ast.Node)                                                                    // Required
	GetEnumMemberValueCache            func(node *ast.Node) EvaluatorResult                                                    // Required
	SetEnumMemberValueCache            func(node *ast.Node, result EvaluatorResult)                                            // Required
	SymbolToString                     func(s *ast.Symbol) string                                                              // Required
	Error                              func(location *ast.Node, message *diagnostics.Message, args ...any) *ast.Diagnostic     // Optional
	GetGlobalSymbol                    func(name string, meaning ast.SymbolFlags, diagnostic *diagnostics.Message) *ast.Symbol // Optional
	IsBlockScopedNameDeclaredBeforeUse func(declaration *ast.Node, usage *ast.Node) bool                                       // Optional
	GetCombinedNodeFlags               func(node *ast.Node) ast.NodeFlags                                                      // Optional
	OnConstantEnumMemberValueComputed  func(member *ast.Node, result EvaluatorResult)                                          // Optional
}

func (c *ConstantEvaluator) Evaluate(expr *ast.Node, location *ast.Node) EvaluatorResult {
	if c.evaluate == nil {
		c.evaluate = newEvaluator(c.evaluateEntity)
	}
	return c.evaluate(expr, location)
}

func (c *ConstantEvaluator) evaluateEntity(expr *ast.Node, location *ast.Node) EvaluatorResult {
	switch expr.Kind {
	case ast.KindIdentifier, ast.KindPropertyAccessExpression:
		symbol := c.NameResolver.ResolveEntityName(expr, ast.SymbolFlagsValue, true /*ignoreErrors*/, false, nil)
		if symbol == nil {
			return EvaluatorResult{nil, false, false, false}
		}
		if expr.Kind == ast.KindIdentifier {
			if isInfinityOrNaNString(expr.Text()) && (symbol == c.getGlobalSymbol(expr.Text(), ast.SymbolFlagsValue, nil /*diagnostic*/)) {
				// Technically we resolved a global lib file here, but the decision to treat this as numeric
				// is more predicated on the fact that the single-file resolution *didn't* resolve to a
				// different meaning of `Infinity` or `NaN`. Transpilers handle this no problem.
				return EvaluatorResult{jsnum.FromString(expr.Text()), false, false, false}
			}
		}
		if symbol.Flags&ast.SymbolFlagsEnumMember != 0 {
			if location != nil {
				return c.evaluateEnumMember(expr, symbol, location)
			}
			return c.GetEnumMemberValue(symbol.ValueDeclaration)
		}
		if c.isConstantVariable(symbol) {
			declaration := symbol.ValueDeclaration
			if declaration != nil && declaration.Type() == nil && declaration.Initializer() != nil &&
				(location == nil || declaration != location && c.IsBlockScopedNameDeclaredBeforeUse(declaration, location)) {
				result := c.evaluate(declaration.Initializer(), declaration)
				if location != nil && ast.GetSourceFileOfNode(location) != ast.GetSourceFileOfNode(declaration) {
					return EvaluatorResult{result.Value, false, true, true}
				}
				return EvaluatorResult{result.Value, result.IsSyntacticallyString, result.ResolvedOtherFiles, true /*hasExternalReferences*/}
			}
		}
		return EvaluatorResult{nil, false, false, false}
	case ast.KindElementAccessExpression:
		root := expr.Expression()
		if ast.IsEntityNameExpression(root) && ast.IsStringLiteralLike(expr.AsElementAccessExpression().ArgumentExpression) {
			rootSymbol := c.NameResolver.ResolveEntityName(root, ast.SymbolFlagsValue, true /*ignoreErrors*/, false, nil)
			if rootSymbol != nil && rootSymbol.Flags&ast.SymbolFlagsEnum != 0 {
				name := expr.AsElementAccessExpression().ArgumentExpression.Text()
				member := rootSymbol.Exports[name]
				if member != nil {
					// !!! Debug.assert(ast.GetSourceFileOfNode(member.valueDeclaration) == ast.GetSourceFileOfNode(rootSymbol.valueDeclaration))
					if location != nil {
						return c.evaluateEnumMember(expr, member, location)
					}
					return c.GetEnumMemberValue(member.ValueDeclaration)
				}
			}
		}
		return EvaluatorResult{nil, false, false, false}
	}
	panic("Unhandled case in evaluateEntity")
}

func (c *ConstantEvaluator) ComputeEnumMemberValues(node *ast.Node) {
	if !c.GetEnumDeclarationValuesComputed(node) {
		c.SetEnumDeclarationValuesComputed(node)
		var autoValue jsnum.Number
		var previous *ast.Node
		for _, member := range node.AsEnumDeclaration().Members.Nodes {
			result := c.computeEnumMemberValue(member, autoValue, previous)
			c.SetEnumMemberValueCache(member, result)
			if value, isNumber := result.Value.(jsnum.Number); isNumber {
				autoValue = value + 1
			} else {
				autoValue = jsnum.NaN()
			}
			previous = member
		}
	}
}

func (c *ConstantEvaluator) GetEnumMemberValue(node *ast.Node) EvaluatorResult {
	if node.Parent != nil {
		c.ComputeEnumMemberValues(node.Parent)
	}
	return c.GetEnumMemberValueCache(node)
}

func (c *ConstantEvaluator) evaluateEnumMember(expr *ast.Node, symbol *ast.Symbol, location *ast.Node) EvaluatorResult {
	declaration := symbol.ValueDeclaration
	if declaration == nil || declaration == location {
		c.error(expr, diagnostics.Property_0_is_used_before_being_assigned, c.SymbolToString(symbol))
		return EvaluatorResult{nil, false, false, false}
	}
	if !c.IsBlockScopedNameDeclaredBeforeUse(declaration, location) {
		c.error(expr, diagnostics.A_member_initializer_in_a_enum_declaration_cannot_reference_members_declared_after_it_including_members_defined_in_other_enums)
		return EvaluatorResult{0.0, false, false, false}
	}
	value := c.GetEnumMemberValue(declaration)
	if location.Parent != declaration.Parent {
		return EvaluatorResult{value.Value, value.IsSyntacticallyString, value.ResolvedOtherFiles, true /*hasExternalReferences*/}
	}
	return value
}

func (c *ConstantEvaluator) error(location *ast.Node, message *diagnostics.Message, args ...any) {
	if c.Error != nil {
		c.Error(location, message, args...)
	}
}

func (c *ConstantEvaluator) computeEnumMemberValue(member *ast.Node, autoValue jsnum.Number, previous *ast.Node) EvaluatorResult {
	if ast.IsComputedNonLiteralName(member.Name()) {
		c.error(member.Name(), diagnostics.Computed_property_names_are_not_allowed_in_enums)
	} else {
		text := member.Name().Text()
		if jsnum.IsNumericLiteralName(text) && !isInfinityOrNaNString(text) {
			c.error(member.Name(), diagnostics.An_enum_member_cannot_have_a_numeric_name)
		}
	}
	if member.Initializer() != nil {
		return c.computeConstantEnumMemberValue(member)
	}
	// In ambient non-const numeric enum declarations, enum members without initializers are
	// considered computed members (as opposed to having auto-incremented values).
	if member.Parent != nil && member.Parent.Flags&ast.NodeFlagsAmbient != 0 && !ast.IsEnumConst(member.Parent) {
		return EvaluatorResult{nil, false, false, false}
	}
	// If the member declaration specifies no value, the member is considered a constant enum member.
	// If the member is the first member in the enum declaration, it is assigned the value zero.
	// Otherwise, it is assigned the value of the immediately preceding member plus one, and an error
	// occurs if the immediately preceding member is not a constant enum member.
	if autoValue.IsNaN() {
		c.error(member.Name(), diagnostics.Enum_member_must_have_initializer)
		return EvaluatorResult{nil, false, false, false}
	}
	if c.CompilerOptions.GetIsolatedModules() && previous != nil && previous.AsEnumMember().Initializer != nil {
		prevValue := c.GetEnumMemberValue(previous)
		_, prevIsNum := prevValue.Value.(jsnum.Number)
		if !prevIsNum || prevValue.ResolvedOtherFiles {
			c.error(member.Name(), diagnostics.Enum_member_following_a_non_literal_numeric_member_must_have_an_initializer_when_isolatedModules_is_enabled)
		}
	}
	return EvaluatorResult{autoValue, false, false, false}
}

func (c *ConstantEvaluator) computeConstantEnumMemberValue(member *ast.Node) EvaluatorResult {
	result := c.Evaluate(member.Initializer(), member)
	if c.OnConstantEnumMemberValueComputed != nil {
		c.OnConstantEnumMemberValueComputed(member, result)
	}
	return result
}

func (c *ConstantEvaluator) getGlobalSymbol(name string, meaning ast.SymbolFlags, diagnostic *diagnostics.Message) *ast.Symbol {
	if c.GetGlobalSymbol != nil {
		return c.GetGlobalSymbol(name, meaning, diagnostic)
	}
	return nil
}

func (c *ConstantEvaluator) isConstantVariable(symbol *ast.Symbol) bool {
	return symbol.Flags&ast.SymbolFlagsVariable != 0 && (c.getDeclarationNodeFlagsFromSymbol(symbol)&ast.NodeFlagsConstant) != 0
}

func (c *ConstantEvaluator) getDeclarationNodeFlagsFromSymbol(s *ast.Symbol) ast.NodeFlags {
	if s.ValueDeclaration != nil {
		return c.getCombinedNodeFlags(s.ValueDeclaration)
	}
	return ast.NodeFlagsNone
}

func (c *ConstantEvaluator) getCombinedNodeFlags(node *ast.Node) ast.NodeFlags {
	if c.GetCombinedNodeFlags != nil {
		return c.GetCombinedNodeFlags(node)
	}
	return ast.GetCombinedNodeFlags(node)
}

func isInfinityOrNaNString(name string) bool {
	return name == "Infinity" || name == "-Infinity" || name == "NaN"
}
