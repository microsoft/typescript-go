// Package diagnostics contains generated localizable diagnostic messages.
package diagnostics

import (
	"fmt"
	"regexp"
	"strconv"

	"golang.org/x/text/language"
)

//go:generate go run generate.go -output ./diagnostics_generated.go
//go:generate go tool golang.org/x/tools/cmd/stringer -type=Category -output=stringer_generated.go
//go:generate go tool mvdan.cc/gofumpt -w diagnostics_generated.go stringer_generated.go

type Category int32

const (
	CategoryWarning Category = iota
	CategoryError
	CategorySuggestion
	CategoryMessage
)

func (category Category) Name() string {
	switch category {
	case CategoryWarning:
		return "warning"
	case CategoryError:
		return "error"
	case CategorySuggestion:
		return "suggestion"
	case CategoryMessage:
		return "message"
	}
	panic("Unhandled diagnostic category")
}

type Key string

type Message struct {
	code                         int32
	category                     Category
	key                          Key
	text                         string
	reportsUnnecessary           bool
	elidedInCompatibilityPyramid bool
	reportsDeprecated            bool
}

func (m *Message) Code() int32                        { return m.code }
func (m *Message) Category() Category                 { return m.category }
func (m *Message) Key() Key                           { return m.key }
func (m *Message) ReportsUnnecessary() bool           { return m.reportsUnnecessary }
func (m *Message) ElidedInCompatibilityPyramid() bool { return m.elidedInCompatibilityPyramid }
func (m *Message) ReportsDeprecated() bool            { return m.reportsDeprecated }

// For debugging only.
func (m *Message) String() string {
	return m.text
}

func (m *Message) Localize(locale language.Tag, args ...any) string {
	return Localize(locale, m, "", StringifyArgs(args)...)
}

func (m *Message) LocalizeStringArgs(locale language.Tag, args ...string) string {
	return Localize(locale, m, "", args...)
}

func Localize(locale language.Tag, message *Message, key Key, args ...string) string {
	if message == nil {
		message = keyToMessage(key)
	}

	if message != nil {
		// !!! localize
		return Format(message.text, args)
	}

	panic("Unknown diagnostic message: " + string(key))
}

var placeholderRegexp = regexp.MustCompile(`{(\d+)}`)

func Format(text string, args []string) string {
	if len(args) == 0 {
		return text
	}

	return placeholderRegexp.ReplaceAllStringFunc(text, func(match string) string {
		index, err := strconv.ParseInt(match[1:len(match)-1], 10, 0)
		if err != nil || int(index) >= len(args) {
			panic("Invalid formatting placeholder")
		}
		return args[int(index)]
	})
}

func StringifyArgs(args []any) []string {
	if len(args) == 0 {
		return nil
	}

	result := make([]string, len(args))
	for i, arg := range args {
		if s, ok := arg.(string); ok {
			result[i] = s
		} else {
			result[i] = fmt.Sprintf("%v", arg)
		}
	}
	return result
}
