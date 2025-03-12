package ast

type NamespaceFlags uint32

const (
	NamespaceFlagsNone               NamespaceFlags = 0
	NamespaceFlagsNestedNamespace    NamespaceFlags = 1 << 1 // Namespace declaration
	NamespaceFlagsNamespace          NamespaceFlags = 1 << 2 // Namespace declaration
	NamespaceFlagsGlobalAugmentation NamespaceFlags = 1 << 3 // Set if module declaration is an augmentation for the global scope
)
