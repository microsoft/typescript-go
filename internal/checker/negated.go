package checker

import (
	"slices"

	"github.com/microsoft/typescript-go/internal/core"
)

// newNegatedType creates a new negated type 'not baseType'.
func (c *Checker) newNegatedType(baseType *Type) *Type {
	data := &NegatedType{}
	data.baseType = baseType
	return c.newType(TypeFlagsNegated, ObjectFlagsNone, data)
}

// getNegatedType constructs the type 'not T', applying the negation identities:
//
//	not any      => any
//	not unknown  => never
//	not never    => unknown
//	not not T    => T
//	not (A | B)  => not A & not B   (De Morgan)
//
// All other types are wrapped in a (cached) NegatedType.
func (c *Checker) getNegatedType(t *Type) *Type {
	switch {
	case t.flags&TypeFlagsAny != 0:
		return t
	case t.flags&TypeFlagsUnknown != 0:
		return c.neverType
	case t.flags&TypeFlagsNever != 0:
		return c.unknownType
	case t.flags&TypeFlagsNegated != 0:
		return t.AsNegatedType().baseType
	case t.flags&TypeFlagsUnion != 0:
		return c.getIntersectionType(core.Map(t.Types(), c.getNegatedType))
	}
	if cached := c.negatedTypes[t.id]; cached != nil {
		return cached
	}
	result := c.newNegatedType(t)
	c.negatedTypes[t.id] = result
	return result
}

// checkForUnsatisfiedNegatedType returns true if the intersection in typeSet is empty (never)
// because some non-negated member is a subtype of the union of the negated members' base types.
// For example, in '"w" & not string' the non-negated member '"w"' is a subtype of 'string'
// (the base type of 'not string'), so the intersection reduces to never.
func (c *Checker) checkForUnsatisfiedNegatedType(typeSet []*Type) bool {
	nonNegatedSet := core.Filter(typeSet, func(t *Type) bool { return t.flags&TypeFlagsNegated == 0 })
	if len(nonNegatedSet) == 0 {
		return false
	}
	negatedBounds := c.getUnionType(core.Map(core.Filter(typeSet, isNegatedType), func(t *Type) *Type {
		return t.AsNegatedType().baseType
	}))
	for _, nonNegatedType := range nonNegatedSet {
		if c.isTypeSubtypeOf(nonNegatedType, negatedBounds) {
			return true
		}
	}
	return false
}

// removeNegatedSubtypes removes redundant negated members from an intersection. A member 'not X'
// is redundant when the combined non-negated part of the intersection is already a subtype of
// 'not X' (i.e. it is disjoint from X). For example, in 'false & not true' the non-negated part
// 'false' is a subtype of 'not true', so 'not true' is dropped, leaving just 'false'.
func (c *Checker) removeNegatedSubtypes(types []*Type) []*Type {
	if len(types) == 0 {
		return types
	}
	nonNegatedBounds := core.Filter(types, func(t *Type) bool { return t.flags&TypeFlagsNegated == 0 })
	if len(nonNegatedBounds) == 0 {
		return types
	}
	nonNegativePart := c.getIntersectionType(nonNegatedBounds)
	for i := len(types) - 1; i >= 0; i-- {
		if types[i].flags&TypeFlagsNegated != 0 && c.isTypeSubtypeOf(nonNegativePart, types[i]) {
			types = slices.Delete(types, i, i+1)
		}
	}
	return types
}
