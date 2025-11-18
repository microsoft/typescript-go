package locale

import (
	"context"

	"golang.org/x/text/language"
)

type contextKey int

type Locale language.Tag

var Default Locale

func WithLocale(ctx context.Context, locale Locale) context.Context {
	return context.WithValue(ctx, contextKey(0), locale)
}

func FromContext(ctx context.Context) Locale {
	locale, _ := ctx.Value(contextKey(0)).(Locale)
	return locale
}

func Parse(localeStr string) Locale {
	// Parse gracefully fails.
	locale, _ := language.Parse(localeStr)
	return Locale(locale)
}
