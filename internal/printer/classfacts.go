package printer

type ClassFacts int32

const (
	ClassFactsClassWasDecorated = 1 << iota
	ClassFactsNeedsClassConstructorReference
	ClassFactsNeedsClassSuperReference
	ClassFactsNeedsSubstitutionForThisInClassStaticField
	ClassFactsWillHoistInitializersToConstructor
	ClassFactsNone ClassFacts = 0
)
