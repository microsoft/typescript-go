package rangesymboltable

import "testdata/fakeast"

func RangeOverSymbolTable(table fakeast.SymbolTable) {
	for range table {
	}
}

func RangeOverSymbolTableKeyValue(table fakeast.SymbolTable) {
	for k, v := range table {
		_ = k
		_ = v
	}
}

func RangeOverSymbolTableKey(table fakeast.SymbolTable) {
	for k := range table {
		_ = k
	}
}

func RangeOverRegularMap(m map[string]int) {
	for range m {
	}
}

func RangeOverSymbolTablePointer(table *fakeast.SymbolTable) {
	// Dereferencing and ranging over a SymbolTable should also be caught.
	for range *table {
	}
}
