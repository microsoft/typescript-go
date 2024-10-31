package compiler

import "github.com/microsoft/typescript-go/internal/compiler/diagnostics"

type IntersectionState uint32

const (
	IntersectionStateNone   IntersectionState = 0
	IntersectionStateSource IntersectionState = 1 << 0 // Source type is a constituent of an outer intersection
	IntersectionStateTarget IntersectionState = 1 << 1 // Target type is a constituent of an outer intersection
)

type RecursionFlags uint32

const (
	RecursionFlagsNone   RecursionFlags = 0
	RecursionFlagsSource RecursionFlags = 1 << 0
	RecursionFlagsTarget RecursionFlags = 1 << 1
	RecursionFlagsBoth                  = RecursionFlagsSource | RecursionFlagsTarget
)

type ExpandingFlags uint8

const (
	ExpandingFlagsNone   ExpandingFlags = 0
	ExpandingFlagsSource ExpandingFlags = 1
	ExpandingFlagsTarget ExpandingFlags = 1 << 1
	ExpandingFlagsBoth                  = ExpandingFlagsSource | ExpandingFlagsTarget
)

type RelationComparisonResult uint32

const (
	RelationComparisonResultNone                RelationComparisonResult = 0
	RelationComparisonResultSucceeded           RelationComparisonResult = 1 << 0
	RelationComparisonResultFailed              RelationComparisonResult = 1 << 1
	RelationComparisonResultReportsUnmeasurable RelationComparisonResult = 1 << 3
	RelationComparisonResultReportsUnreliable   RelationComparisonResult = 1 << 4
	RelationComparisonResultComplexityOverflow  RelationComparisonResult = 1 << 5
	RelationComparisonResultStackDepthOverflow  RelationComparisonResult = 1 << 6
	RelationComparisonResultReportsMask                                  = RelationComparisonResultReportsUnmeasurable | RelationComparisonResultReportsUnreliable
	RelationComparisonResultOverflow                                     = RelationComparisonResultComplexityOverflow | RelationComparisonResultStackDepthOverflow
)

type DiagnosticAndArguments struct {
	message   *diagnostics.Message
	arguments []any
}

type ErrorOutputContainer struct {
	errors      []*Diagnostic
	skipLogging bool
}

type ErrorReporter func(message *diagnostics.Message, args ...any)

type Relation struct {
	results map[string]RelationComparisonResult
}

func (r *Relation) get(key string) RelationComparisonResult {
	return r.results[key]
}

func (r *Relation) set(key string, result RelationComparisonResult) {
	if r.results == nil {
		r.results = make(map[string]RelationComparisonResult)
	}
	r.results[key] = result
}

func (c *Checker) isTypeAssignableTo(source *Type, target *Type) bool {
	return c.isTypeRelatedTo(source, target, c.assignableRelation)
}

func (c *Checker) isTypeStrictSubtypeOf(source *Type, target *Type) bool {
	return c.isTypeRelatedTo(source, target, c.strictSubtypeRelation)
}

func (c *Checker) isTypeRelatedTo(source *Type, target *Type, relation *Relation) bool {
	if isFreshLiteralType(source) {
		source = source.AsLiteralType().regularType
	}
	if isFreshLiteralType(target) {
		target = target.AsLiteralType().regularType
	}
	if source == target {
		return true
	}
	if relation != c.identityRelation {
		if relation == c.comparableRelation && target.flags&TypeFlagsNever == 0 && c.isSimpleTypeRelatedTo(target, source, relation, nil) || c.isSimpleTypeRelatedTo(source, target, relation, nil) {
			return true
		}
	} else if !((source.flags|target.flags)&(TypeFlagsUnionOrIntersection|TypeFlagsIndexedAccess|TypeFlagsConditional|TypeFlagsSubstitution) != 0) {
		// We have excluded types that may simplify to other forms, so types must have identical flags
		if source.flags != target.flags {
			return false
		}
		if source.flags&TypeFlagsSingleton != 0 {
			return true
		}
	}
	if source.flags&TypeFlagsObject != 0 && target.flags&TypeFlagsObject != 0 {
		related := relation.get(getRelationKey(source, target, IntersectionStateNone, relation == c.identityRelation, false))
		if related != RelationComparisonResultNone {
			return related&RelationComparisonResultSucceeded != 0
		}
	}
	if source.flags&TypeFlagsStructuredOrInstantiable != 0 || target.flags&TypeFlagsStructuredOrInstantiable != 0 {
		return c.checkTypeRelatedTo(source, target, relation, nil /*errorNode*/)
	}
	return false
}

func (c *Checker) isSimpleTypeRelatedTo(source *Type, target *Type, relation *Relation, errorReporter ErrorReporter) bool {
	s := source.flags
	t := target.flags
	if t&TypeFlagsAny != 0 || s&TypeFlagsNever != 0 || source == c.wildcardType {
		return true
	}
	if t&TypeFlagsUnknown != 0 && !(relation == c.strictSubtypeRelation && s&TypeFlagsAny != 0) {
		return true
	}
	if t&TypeFlagsNever != 0 {
		return false
	}
	if s&TypeFlagsStringLike != 0 && t&TypeFlagsString != 0 {
		return true
	}
	if s&TypeFlagsStringLiteral != 0 && s&TypeFlagsEnumLiteral != 0 && t&TypeFlagsStringLiteral != 0 && t&TypeFlagsEnumLiteral == 0 && source.AsLiteralType().value == target.AsLiteralType().value {
		return true
	}
	if s&TypeFlagsNumberLike != 0 && t&TypeFlagsNumber != 0 {
		return true
	}
	if s&TypeFlagsNumberLiteral != 0 && s&TypeFlagsEnumLiteral != 0 && t&TypeFlagsNumberLiteral != 0 && t&TypeFlagsEnumLiteral == 0 && source.AsLiteralType().value == target.AsLiteralType().value {
		return true
	}
	if s&TypeFlagsBigIntLike != 0 && t&TypeFlagsBigInt != 0 {
		return true
	}
	if s&TypeFlagsBooleanLike != 0 && t&TypeFlagsBoolean != 0 {
		return true
	}
	if s&TypeFlagsESSymbolLike != 0 && t&TypeFlagsESSymbol != 0 {
		return true
	}
	if s&TypeFlagsEnum != 0 && t&TypeFlagsEnum != 0 && source.symbol.name == target.symbol.name && c.isEnumTypeRelatedTo(source.symbol, target.symbol, errorReporter) {
		return true
	}
	if s&TypeFlagsEnumLiteral != 0 && t&TypeFlagsEnumLiteral != 0 {
		if s&TypeFlagsUnion != 0 && t&TypeFlagsUnion != 0 && c.isEnumTypeRelatedTo(source.symbol, target.symbol, errorReporter) {
			return true
		}
		if s&TypeFlagsLiteral != 0 && t&TypeFlagsLiteral != 0 && source.AsLiteralType().value == target.AsLiteralType().value && c.isEnumTypeRelatedTo(source.symbol, target.symbol, errorReporter) {
			return true
		}
	}
	// In non-strictNullChecks mode, `undefined` and `null` are assignable to anything except `never`.
	// Since unions and intersections may reduce to `never`, we exclude them here.
	if s&TypeFlagsUndefined != 0 && (!c.strictNullChecks && t&TypeFlagsUnionOrIntersection == 0 || t&(TypeFlagsUndefined|TypeFlagsVoid) != 0) {
		return true
	}
	if s&TypeFlagsNull != 0 && (!c.strictNullChecks && t&TypeFlagsUnionOrIntersection == 0 || t&TypeFlagsNull != 0) {
		return true
	}
	if s&TypeFlagsObject != 0 && t&TypeFlagsNonPrimitive != 0 && !(relation == c.strictSubtypeRelation && c.isEmptyAnonymousObjectType(source) && source.objectFlags&ObjectFlagsFreshLiteral == 0) {
		return true
	}
	if relation == c.assignableRelation || relation == c.comparableRelation {
		if s&TypeFlagsAny != 0 {
			return true
		}
		// Type number is assignable to any computed numeric enum type or any numeric enum literal type, and
		// a numeric literal type is assignable any computed numeric enum type or any numeric enum literal type
		// with a matching value. These rules exist such that enums can be used for bit-flag purposes.
		if s&TypeFlagsNumber != 0 && (t&TypeFlagsEnum != 0 || t&TypeFlagsNumberLiteral != 0 && t&TypeFlagsEnumLiteral != 0) {
			return true
		}
		if s&TypeFlagsNumberLiteral != 0 && s&TypeFlagsEnumLiteral == 0 && (t&TypeFlagsEnum != 0 || t&TypeFlagsNumberLiteral != 0 && t&TypeFlagsEnumLiteral != 0 && source.AsLiteralType().value == target.AsLiteralType().value) {
			return true
		}
		// Anything is assignable to a union containing undefined, null, and {}
		if c.isUnknownLikeUnionType(target) {
			return true
		}
	}
	return false
}

func (c *Checker) isEnumTypeRelatedTo(source *Symbol, target *Symbol, errorReporter ErrorReporter) bool {
	return source == target // !!!
}

func (c *Checker) checkTypeRelatedTo(source *Type, target *Type, relation *Relation, errorNode *Node) bool {
	return c.checkTypeRelatedToEx(source, target, relation, errorNode, nil, nil, nil)
}

func (c *Checker) checkTypeRelatedToEx(
	source *Type,
	target *Type,
	relation *Relation,
	errorNode *Node,
	headMessage *diagnostics.Message,
	containingMessageChain func() *MessageChain,
	errorOutputContainer *ErrorOutputContainer,
) bool {
	r := Relater{c: c, relation: relation}
	result := r.isRelatedToEx(source, target, RecursionFlagsBoth, errorNode != nil /*reportErrors*/, headMessage, IntersectionStateNone)
	if len(r.incompatibleStack) != 0 {
		r.reportIncompatibleStack()
	}
	if r.overflow {
		// Record this relation as having failed such that we don't attempt the overflowing operation again.
		id := getRelationKey(source, target, IntersectionStateNone, relation == c.identityRelation, false /*ignoreConstraints*/)
		relation.set(id, RelationComparisonResultFailed|ifElse(r.relationCount <= 0, RelationComparisonResultComplexityOverflow, RelationComparisonResultStackDepthOverflow))
		message := ifElse(r.relationCount <= 0, diagnostics.Excessive_complexity_comparing_types_0_and_1, diagnostics.Excessive_stack_depth_comparing_types_0_and_1)
		if errorNode == nil {
			errorNode = c.currentNode
		}
		diag := c.error(errorNode, message, c.typeToString(source), c.typeToString(target))
		if errorOutputContainer != nil {
			errorOutputContainer.errors = append(errorOutputContainer.errors, diag)
		}
	} else if r.errorInfo != nil {
		if containingMessageChain != nil {
			chain := containingMessageChain()
			if chain != nil {
				concatenateDiagnosticMessageChains(chain, r.errorInfo)
				r.errorInfo = chain
			}
		}
		// !!!
		// var relatedInformation []*Diagnostic
		// // Check if we should issue an extra diagnostic to produce a quickfix for a slightly incorrect import statement
		// if headMessage != nil && errorNode != nil && result == TernaryFalse && source.symbol != nil {
		// 	links := c.getSymbolLinks(source.symbol)
		// 	if links.originatingImport && !isImportCall(links.originatingImport) {
		// 		helpfulRetry := c.checkTypeRelatedTo(c.getTypeOfSymbol(links.target), target, relation /*errorNode*/, nil)
		// 		if helpfulRetry {
		// 			// Likely an incorrect import. Issue a helpful diagnostic to produce a quickfix to change the import
		// 			diag := createDiagnosticForNode(links.originatingImport, Diagnostics.Type_originates_at_this_import_A_namespace_style_import_cannot_be_called_or_constructed_and_will_cause_a_failure_at_runtime_Consider_using_a_default_import_or_import_require_here_instead)
		// 			relatedInformation = append(relatedInformation, diag)
		// 			// Cause the error to appear with the error that triggered it
		// 		}
		// 	}
		// }
		diag := NewDiagnosticForNodeFromMessageChain(errorNode, r.errorInfo).setRelatedInfo(r.relatedInfo)
		if errorOutputContainer != nil {
			errorOutputContainer.errors = append(errorOutputContainer.errors, diag)
		}
		if errorOutputContainer == nil || !errorOutputContainer.skipLogging {
			c.diagnostics.add(diag)
		}
	}
	// !!!
	// if errorNode != nil && errorOutputContainer != nil && errorOutputContainer.skipLogging && result == TernaryFalse {
	// 	Debug.assert(!!errorOutputContainer.errors, "missed opportunity to interact with error.")
	// }
	return result != TernaryFalse
}

// A type is 'weak' if it is an object type with at least one optional property
// and no required properties, call/construct signatures or index signatures
func (c *Checker) isWeakType(t *Type) bool {
	if t.flags&TypeFlagsObject != 0 {
		resolved := c.resolveStructuredTypeMembers(t)
		return len(resolved.signatures) == 0 && len(resolved.indexInfos) == 0 && len(resolved.properties) > 0 && every(resolved.properties, func(p *Symbol) bool {
			return p.flags&SymbolFlagsOptional != 0
		})
	}
	if t.flags&TypeFlagsSubstitution != 0 {
		return c.isWeakType(t.AsSubstitutionType().baseType)
	}
	if t.flags&TypeFlagsIntersection != 0 {
		return every(t.AsIntersectionType().types, c.isWeakType)
	}
	return false
}

func (c *Checker) hasCommonProperties(source *Type, target *Type, isComparingJsxAttributes bool) bool {
	return false // !!!
}

type Relater struct {
	c                     *Checker
	relation              *Relation
	errorInfo             *MessageChain
	relatedInfo           []*Diagnostic
	maybeKeys             []string
	maybeKeysSet          set[string]
	sourceStack           []*Type
	targetStack           []*Type
	maybeCount            int
	sourceDepth           int
	targetDepth           int
	expandingFlags        ExpandingFlags
	overflow              bool
	overrideNextErrorInfo int
	skipParentCounter     int // How many `reportRelationError` calls should be skipped in the elaboration pyramid
	lastSkippedInfo       [2]*Type
	incompatibleStack     []DiagnosticAndArguments
	relationCount         int
}

func (r *Relater) isRelatedTo(originalSource *Type, originalTarget *Type, recursionFlags RecursionFlags, reportErrors bool) Ternary {
	return r.isRelatedToEx(originalSource, originalTarget, recursionFlags, reportErrors, nil, IntersectionStateNone)
}

func (r *Relater) isRelatedToEx(originalSource *Type, originalTarget *Type, recursionFlags RecursionFlags, reportErrors bool, headMessage *diagnostics.Message, intersectionState IntersectionState) Ternary {
	if originalSource == originalTarget {
		return TernaryTrue
	}
	// Before normalization: if `source` is type an object type, and `target` is primitive,
	// skip all the checks we don't need and just return `isSimpleTypeRelatedTo` result
	if originalSource.flags&TypeFlagsObject != 0 && originalTarget.flags&TypeFlagsPrimitive != 0 {
		if r.relation == r.c.comparableRelation && originalTarget.flags&TypeFlagsNever == 0 && r.c.isSimpleTypeRelatedTo(originalTarget, originalSource, r.relation, nil) ||
			r.c.isSimpleTypeRelatedTo(originalSource, originalTarget, r.relation, ifElse(reportErrors, r.reportError, nil)) {
			return TernaryTrue
		}
		if reportErrors {
			r.reportErrorResults(originalSource, originalTarget, originalSource, originalTarget, headMessage)
		}
		return TernaryFalse
	}
	// Normalize the source and target types: Turn fresh literal types into regular literal types,
	// turn deferred type references into regular type references, simplify indexed access and
	// conditional types, and resolve substitution types to either the substitution (on the source
	// side) or the type variable (on the target side).
	source := r.c.getNormalizedType(originalSource, false /*writing*/)
	target := r.c.getNormalizedType(originalTarget, true /*writing*/)
	if source == target {
		return TernaryTrue
	}
	if r.relation == r.c.identityRelation {
		if source.flags != target.flags {
			return TernaryFalse
		}
		if source.flags&TypeFlagsSingleton != 0 {
			return TernaryTrue
		}
		// !!! traceUnionsOrIntersectionsTooLarge(source, target)
		return r.recursiveTypeRelatedTo(source, target /*reportErrors*/, false, IntersectionStateNone, recursionFlags)
	}
	// We fastpath comparing a type parameter to exactly its constraint, as this is _super_ common,
	// and otherwise, for type parameters in large unions, causes us to need to compare the union to itself,
	// as we break down the _target_ union first, _then_ get the source constraint - so for every
	// member of the target, we attempt to find a match in the source. This avoids that in cases where
	// the target is exactly the constraint.
	if source.flags&TypeFlagsTypeParameter != 0 && r.c.getConstraintOfType(source) == target {
		return TernaryTrue
	}
	// See if we're relating a definitely non-nullable type to a union that includes null and/or undefined
	// plus a single non-nullable type. If so, remove null and/or undefined from the target type.
	if source.flags&TypeFlagsDefinitelyNonNullable != 0 && target.flags&TypeFlagsUnion != 0 {
		types := target.AsUnionType().types
		var candidate *Type
		switch {
		case len(types) == 2 && types[0].flags&TypeFlagsNullable != 0:
			candidate = types[1]
		case len(types) == 3 && types[0].flags&TypeFlagsNullable != 0 && types[1].flags&TypeFlagsNullable != 0:
			candidate = types[2]
		}
		if candidate != nil && candidate.flags&TypeFlagsNullable == 0 {
			target = r.c.getNormalizedType(candidate /*writing*/, true)
			if source == target {
				return TernaryTrue
			}
		}
	}
	if r.relation == r.c.comparableRelation && target.flags&TypeFlagsNever == 0 && r.c.isSimpleTypeRelatedTo(target, source, r.relation, nil) ||
		r.c.isSimpleTypeRelatedTo(source, target, r.relation, ifElse(reportErrors, r.reportError, nil)) {
		return TernaryTrue
	}
	if source.flags&TypeFlagsStructuredOrInstantiable != 0 || target.flags&TypeFlagsStructuredOrInstantiable != 0 {
		isPerformingExcessPropertyChecks := intersectionState&IntersectionStateTarget == 0 && isObjectLiteralType(source) && source.objectFlags&ObjectFlagsFreshLiteral != 0
		if isPerformingExcessPropertyChecks {
			if r.hasExcessProperties(source, target, reportErrors) {
				if reportErrors {
					r.reportRelationError(headMessage, source, ifElse(originalTarget.alias != nil, originalTarget, target))
				}
				return TernaryFalse
			}
		}
		isPerformingCommonPropertyChecks := (r.relation != r.c.comparableRelation || isUnitType(source)) &&
			intersectionState&IntersectionStateTarget == 0 &&
			source.flags&(TypeFlagsPrimitive|TypeFlagsObject|TypeFlagsIntersection) != 0 && source != r.c.globalObjectType &&
			target.flags&(TypeFlagsObject|TypeFlagsIntersection) != 0 && r.c.isWeakType(target) && (len(r.c.getPropertiesOfType(source)) > 0 || r.c.typeHasCallOrConstructSignatures(source))
		isComparingJsxAttributes := source.objectFlags&ObjectFlagsJsxAttributes != 0
		if isPerformingCommonPropertyChecks && !r.c.hasCommonProperties(source, target, isComparingJsxAttributes) {
			if reportErrors {
				sourceString := r.c.typeToString(ifElse(originalSource.alias != nil, originalSource, source))
				targetString := r.c.typeToString(ifElse(originalTarget.alias != nil, originalTarget, target))
				calls := r.c.getSignaturesOfType(source, SignatureKindCall)
				constructs := r.c.getSignaturesOfType(source, SignatureKindConstruct)
				if len(calls) > 0 && r.isRelatedTo(r.c.getReturnTypeOfSignature(calls[0]), target, RecursionFlagsSource, false /*reportErrors*/) != TernaryFalse ||
					len(constructs) > 0 && r.isRelatedTo(r.c.getReturnTypeOfSignature(constructs[0]), target, RecursionFlagsSource, false /*reportErrors*/) != TernaryFalse {
					r.reportError(diagnostics.Value_of_type_0_has_no_properties_in_common_with_type_1_Did_you_mean_to_call_it, sourceString, targetString)
				} else {
					r.reportError(diagnostics.Type_0_has_no_properties_in_common_with_type_1, sourceString, targetString)
				}
			}
			return TernaryFalse
		}
		// !!! traceUnionsOrIntersectionsTooLarge(source, target)
		skipCaching := source.flags&TypeFlagsUnion != 0 && len(source.AsUnionType().types) < 4 && target.flags&TypeFlagsUnion == 0 ||
			target.flags&TypeFlagsUnion != 0 && len(target.AsUnionType().types) < 4 && source.flags&TypeFlagsStructuredOrInstantiable == 0
		var result Ternary
		if skipCaching {
			result = r.unionOrIntersectionRelatedTo(source, target, reportErrors, intersectionState)
		} else {
			result = r.recursiveTypeRelatedTo(source, target, reportErrors, intersectionState, recursionFlags)
		}
		if result != TernaryFalse {
			return result
		}
	}
	if reportErrors {
		r.reportErrorResults(originalSource, originalTarget, source, target, headMessage)
	}
	return TernaryFalse
}

func (r *Relater) hasExcessProperties(source *Type, target *Type, reportErrors bool) bool {
	return false // !!!
}

func (r *Relater) unionOrIntersectionRelatedTo(source *Type, target *Type, reportErrors bool, intersectionState IntersectionState) Ternary {
	return TernaryFalse // !!!
}

// Determine if possibly recursive types are related. First, check if the result is already available in the global cache.
// Second, check if we have already started a comparison of the given two types in which case we assume the result to be true.
// Third, check if both types are part of deeply nested chains of generic type instantiations and if so assume the types are
// equal and infinitely expanding. Fourth, if we have reached a depth of 100 nested comparisons, assume we have runaway recursion
// and issue an error. Otherwise, actually compare the structure of the two types.
func (r *Relater) recursiveTypeRelatedTo(source *Type, target *Type, reportErrors bool, intersectionState IntersectionState, recursionFlags RecursionFlags) Ternary {
	return TernaryFalse // !!!
}

func (r *Relater) reportErrorResults(originalSource *Type, originalTarget *Type, source *Type, target *Type, headMessage *diagnostics.Message) {
	sourceHasBase := r.c.getSingleBaseForNonAugmentingSubtype(originalSource) != nil
	targetHasBase := r.c.getSingleBaseForNonAugmentingSubtype(originalTarget) != nil
	if originalSource.alias != nil || sourceHasBase {
		source = originalSource
	}
	if originalTarget.alias != nil || targetHasBase {
		target = originalTarget
	}
	// !!!
	r.reportRelationError(headMessage, source, target)
}

func (r *Relater) reportRelationError(message *diagnostics.Message, source *Type, target *Type) {
	if len(r.incompatibleStack) != 0 {
		r.reportIncompatibleStack()
	}
	// !!!
	sourceType := r.c.typeToString(source)
	targetType := r.c.typeToString(target)
	if message == nil {
		if r.relation == r.c.comparableRelation {
			message = diagnostics.Type_0_is_not_comparable_to_type_1
		} else if sourceType == targetType {
			message = diagnostics.Type_0_is_not_assignable_to_type_1_Two_different_types_with_this_name_exist_but_they_are_unrelated
		} else {
			message = diagnostics.Type_0_is_not_assignable_to_type_1
		}
	}
	r.reportError(message, sourceType, targetType)
}

func (r *Relater) reportError(message *diagnostics.Message, args ...any) {
	// !!! Debug.assert(!!errorNode)
	if len(r.incompatibleStack) != 0 {
		r.reportIncompatibleStack()
	}
	if message.ElidedInCompatabilityPyramid() {
		return
	}
	if r.skipParentCounter == 0 {
		r.errorInfo = NewMessageChain(message, args...).addMessageChain(r.errorInfo)
	} else {
		r.skipParentCounter--
	}
}

func (r *Relater) reportIncompatibleStack() {
	// !!!
}
