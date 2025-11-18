// Package diagnostics contains generated localizable diagnostic messages.
package diagnostics

import "github.com/microsoft/typescript-go/internal/stringutil"

//go:generate go run generate.go -output ./diagnostics_generated.go
//go:generate go tool mvdan.cc/gofumpt -w diagnostics_generated.go

type Message struct {
	code                         int32
	key                          string
	text                         string
	reportsUnnecessary           bool
	elidedInCompatibilityPyramid bool
	reportsDeprecated            bool
}

func (m *Message) Code() int32                        { return m.code }
func (m *Message) Key() string                        { return m.key }
func (m *Message) Message() string                    { return m.text }
func (m *Message) ReportsUnnecessary() bool           { return m.reportsUnnecessary }
func (m *Message) ElidedInCompatibilityPyramid() bool { return m.elidedInCompatibilityPyramid }
func (m *Message) ReportsDeprecated() bool            { return m.reportsDeprecated }

func (m *Message) Format(args ...any) string {
	text := m.Message()
	if len(args) != 0 {
		text = stringutil.Format(text, args)
	}
	return text
}

func FormatMessage(m *Message, args ...any) *Message {
	result := *m
	result.text = stringutil.Format(m.text, args)
	return &result
}
