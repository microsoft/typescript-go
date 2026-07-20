package checker

import "github.com/microsoft/typescript-go/internal/ast"

// transientSymbol co-allocates a checker-created symbol with its value symbol links, reached via
// symbol.CheckerData. Such symbols are owned by the creating Checker alone and must never be
// handed to another checker; binder symbols use c.valueSymbolLinks (only via the accessors).
type transientSymbol struct {
	symbol ast.Symbol
	data   checkerSymbolData
}

type checkerSymbolData struct {
	valueLinks    ValueSymbolLinks
	hasValueLinks bool
}

func (c *Checker) getValueSymbolLinks(symbol *ast.Symbol) *ValueSymbolLinks {
	if symbol != nil {
		if data, ok := symbol.CheckerData.(*checkerSymbolData); ok {
			data.hasValueLinks = true
			return &data.valueLinks
		}
	}
	return c.valueSymbolLinks.Get(symbol)
}

func (c *Checker) tryGetValueSymbolLinks(symbol *ast.Symbol) *ValueSymbolLinks {
	if symbol != nil {
		if data, ok := symbol.CheckerData.(*checkerSymbolData); ok {
			if data.hasValueLinks {
				return &data.valueLinks
			}
			return nil
		}
	}
	return c.valueSymbolLinks.TryGet(symbol)
}
