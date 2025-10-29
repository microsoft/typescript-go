package unexportedapi

type Foo struct {
	Bar *oops
}

type oops struct {
	v int
}

type Okay struct {
	Sure  int
	Value ***Okay2
}

type Okay2 struct {
	VeryGood struct{}
}

func OkayFunc(v *Okay) *Okay2 {
	if v == nil {
		return nil
	}
	return **v.Value
}

// Test cases for various scenarios

// Exported function with unexported parameter type
func BadFunc(x unexported) {}

// Exported function with unexported return type
func AnotherBadFunc() *unexported {
	return nil
}

// Exported function with unexported type in slice
func SliceFunc(x []unexported) {}

// Exported function with unexported type in map
func MapFunc(x map[string]unexported) {}

// Exported function with unexported type in map key
func MapKeyFunc(x map[unexported]string) {}

// Exported function with unexported type in channel
func ChanFunc(x chan unexported) {}

// Exported type alias to unexported type
type BadAlias = unexported

// Exported type with unexported embedded field
type BadEmbed struct {
	unexported
}

// Unexported type - should not trigger
type okayUnexported struct {
	field unexported
}

// Exported interface with unexported type in method
type BadInterface interface {
	Method(x unexported)
}

// Exported interface with unexported return type
type AnotherBadInterface interface {
	Method() unexported
}

type unexported struct {
	x int
}

// Exported function with multiple return values including unexported
func MultiReturn() (int, unexported, error) {
	return 0, unexported{}, nil
}

// Exported variable with unexported type
var BadVar unexported

// Exported const with unexported type (should not be possible, but let's be safe)
// const BadConst unexported = unexported{} // This won't compile anyway

// Array of unexported type
type BadArray [10]unexported

// Exported function with variadic unexported parameter
func VariadicFunc(args ...unexported) {}

// Exported type with method returning unexported type
type ExportedWithMethod struct{}

func (e ExportedWithMethod) Method() unexported {
	return unexported{}
}

// Exported type with pointer receiver method returning unexported type
func (e *ExportedWithMethod) PointerMethod() *unexported {
	return nil
}

// Generic type with unexported type constraint (Go 1.18+)
type GenericExported[T any] struct {
	Value T
}

// Okay - unexported method on exported type (methods are not part of exported API unless on exported interface)
func (e ExportedWithMethod) unexportedMethod() unexported {
	return unexported{}
}

// Test variables initialized with function calls

// Helper functions for testing
func helperReturnsExported() *Okay2 {
	return &Okay2{}
}

func helperReturnsUnexported() unexported {
	return unexported{}
}

// Okay - exported variable initialized by calling unexported function that returns exported type
var OkayVarFromUnexportedFunc = helperReturnsExported()

// Bad - exported variable initialized by calling exported function that returns unexported type
var BadVarFromFunc = helperReturnsUnexported()

// Okay - exported variable with explicit type (implementation doesn't matter)
var OkayVarExplicitType *Okay2 = helperReturnsExported()

// Bad - exported variable with explicit unexported type
var BadVarExplicitType unexported = helperReturnsUnexported()

// Test type aliases
type (
	ExportedString   string
	unexportedString string
)

// Okay - exported function using exported type alias
func OkayTypeAlias(s ExportedString) {}

// Bad - exported function using unexported type alias
func BadTypeAlias(s unexportedString) {}
