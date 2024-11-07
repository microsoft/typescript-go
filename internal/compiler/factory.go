package compiler

import "slices"

func findSpanEnd[T any](slice []T, test func(value T) bool, start int) int {
	i := start
	for i < len(slice) && test(slice[i]) {
		i++
	}
	return i
}

func (f *NodeFactory) MergeLexicalEnvironment(statements []*Statement, declarations []*Statement) []*Statement {
	if len(declarations) == 0 {
		return statements
	}
	if len(statements) == 0 {
		return declarations
	}

	// When we merge new lexical statements into an existing statement list, we merge them in the following manner:
	//
	// Given:
	//
	// | Left                               | Right                               |
	// |------------------------------------|-------------------------------------|
	// | [standard prologues (left)]        | [standard prologues (right)]        |
	// | [hoisted functions (left)]         | [hoisted functions (right)]         |
	// | [hoisted variables (left)]         | [hoisted variables (right)]         |
	// | [lexical init statements (left)]   | [lexical init statements (right)]   |
	// | [other statements (left)]          |                                     |
	//
	// The resulting statement list will be:
	//
	// | Result                              |
	// |-------------------------------------|
	// | [standard prologues (right)]        |
	// | [standard prologues (left)]         |
	// | [hoisted functions (right)]         |
	// | [hoisted functions (left)]          |
	// | [hoisted variables (right)]         |
	// | [hoisted variables (left)]          |
	// | [lexical init statements (right)]   |
	// | [lexical init statements (left)]    |
	// | [other statements (left)]           |
	//
	// NOTE: It is expected that new lexical init statements must be evaluated before existing lexical init statements,
	// as the prior transformation may depend on the evaluation of the lexical init statements to be in the correct state.

	// find standard prologues on left in the following order: standard directives, hoisted functions, hoisted variables, other custom
	leftStandardPrologueEnd := findSpanEnd(statements, isPrologueDirective, 0)
	leftHoistedFunctionsEnd := findSpanEnd(statements, isHoistedFunction, leftStandardPrologueEnd)
	leftHoistedVariablesEnd := findSpanEnd(statements, isHoistedVariableStatement, leftHoistedFunctionsEnd)

	// find standard prologues on right in the following order: standard directives, hoisted functions, hoisted variables, other custom
	rightStandardPrologueEnd := findSpanEnd(declarations, isPrologueDirective, 0)
	rightHoistedFunctionsEnd := findSpanEnd(declarations, isHoistedFunction, rightStandardPrologueEnd)
	rightHoistedVariablesEnd := findSpanEnd(declarations, isHoistedVariableStatement, rightHoistedFunctionsEnd)
	rightCustomPrologueEnd := findSpanEnd(declarations, isCustomPrologue, rightHoistedVariablesEnd)
	_assert(rightCustomPrologueEnd == len(declarations), "Expected declarations to be valid standard or custom prologues")

	// splice prologues from the right into the left. We do this in reverse order
	// so that we don't need to recompute the index on the left when we insert items.
	left := make([]*Statement, len(statements))
	copy(left, statements)

	// splice other custom prologues from right into left
	if rightCustomPrologueEnd > rightHoistedVariablesEnd {
		left = slices.Insert(left, leftHoistedVariablesEnd, declarations[rightHoistedVariablesEnd:rightCustomPrologueEnd]...)
	}

	// splice hoisted variables from right into left
	if rightHoistedVariablesEnd > rightHoistedFunctionsEnd {
		left = slices.Insert(left, leftHoistedFunctionsEnd, declarations[rightHoistedFunctionsEnd:rightHoistedVariablesEnd]...)
	}

	// splice hoisted functions from right into left
	if rightHoistedFunctionsEnd > rightStandardPrologueEnd {
		left = slices.Insert(left, leftStandardPrologueEnd, declarations[rightStandardPrologueEnd:rightHoistedFunctionsEnd]...)
	}

	// splice standard prologues from right into left (that are not already in left)
	if rightStandardPrologueEnd > 0 {
		if leftStandardPrologueEnd == 0 {
			left = slices.Insert(left, 0, declarations[:rightStandardPrologueEnd]...)
		} else {
			var leftPrologues set[string]
			for i := 0; i < leftStandardPrologueEnd; i++ {
				leftPrologue := statements[i]
				leftPrologues.add(leftPrologue.AsExpressionStatement().expression.Text())
			}
			for i := rightStandardPrologueEnd - 1; i >= 0; i-- {
				rightPrologue := declarations[i]
				if !leftPrologues.has(rightPrologue.AsExpressionStatement().expression.Text()) {
					left = slices.Insert(left, 0, rightPrologue)
				}
			}
		}
	}

	return left
}

func (f *NodeFactory) LiftToBlock(nodes []*Statement) *Statement {
	if len(nodes) == 1 {
		return nodes[0]
	}
	return f.NewBlock(nodes, false /*multiline*/)
}

func (f *NodeFactory) GetGeneratedNameForNode(node *Node) *Node {
	// TODO(rbuckton): To be implemented
	return nil
}

func (f *NodeFactory) CloneNode(node *Node) *Node {
	// TODO(rbuckton): To be implemented
	return nil
}
