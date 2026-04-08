package core

import (
	"context"
)

type key int

const (
	requestIDKey key = iota
	checkerPurposeKey
)

func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey, id)
}

func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(requestIDKey).(string); ok {
		return id
	}
	return ""
}

// CheckerPurpose indicates why a checker is being requested.
type CheckerPurpose int

const (
	// CheckerPurposeQuery is the default purpose for language service operations
	// such as hover, completions, go-to-definition, etc.
	CheckerPurposeQuery CheckerPurpose = iota
	// CheckerPurposeDiagnostics indicates the checker is being used for diagnostics.
	// Diagnostic checkers are dedicated to ensure consistent walk order.
	CheckerPurposeDiagnostics
)

func WithCheckerPurpose(ctx context.Context, purpose CheckerPurpose) context.Context {
	return context.WithValue(ctx, checkerPurposeKey, purpose)
}

func GetCheckerPurpose(ctx context.Context) CheckerPurpose {
	if purpose, ok := ctx.Value(checkerPurposeKey).(CheckerPurpose); ok {
		return purpose
	}
	return CheckerPurposeQuery
}
