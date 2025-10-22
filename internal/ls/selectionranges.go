package ls

import (
	"context"

	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
)

func (l *LanguageService) ProvideSelectionRanges(ctx context.Context, params *lsproto.SelectionRangeParams) (lsproto.SelectionRangeResponse, error) {
	panic("TODO")
}
