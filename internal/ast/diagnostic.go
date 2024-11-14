package ast

import (
	"github.com/microsoft/typescript-go/internal/compiler/diagnostics"
	"github.com/microsoft/typescript-go/internal/core"
)

// Diagnostic

type Diagnostic struct {
	File_               *SourceFile
	Loc_                core.TextRange
	Code_               int32
	Category_           diagnostics.Category
	Message_            string
	MessageChain_       []*MessageChain
	RelatedInformation_ []*Diagnostic
}

func (d *Diagnostic) File() *SourceFile                 { return d.File_ }
func (d *Diagnostic) Pos() int                          { return d.Loc_.Pos() }
func (d *Diagnostic) End() int                          { return d.Loc_.End() }
func (d *Diagnostic) Len() int                          { return d.Loc_.Len() }
func (d *Diagnostic) Loc() core.TextRange               { return d.Loc_ }
func (d *Diagnostic) Code() int32                       { return d.Code_ }
func (d *Diagnostic) Category() diagnostics.Category    { return d.Category_ }
func (d *Diagnostic) Message() string                   { return d.Message_ }
func (d *Diagnostic) MessageChain() []*MessageChain     { return d.MessageChain_ }
func (d *Diagnostic) RelatedInformation() []*Diagnostic { return d.RelatedInformation_ }

func (d *Diagnostic) SetCategory(category diagnostics.Category) { d.Category_ = category }

func (d *Diagnostic) SetMessageChain(messageChain []*MessageChain) *Diagnostic {
	d.MessageChain_ = messageChain
	return d
}

func (d *Diagnostic) AddMessageChain(messageChain *MessageChain) *Diagnostic {
	if messageChain != nil {
		d.MessageChain_ = append(d.MessageChain_, messageChain)
	}
	return d
}

func (d *Diagnostic) SetRelatedInfo(relatedInformation []*Diagnostic) *Diagnostic {
	d.RelatedInformation_ = relatedInformation
	return d
}

func (d *Diagnostic) AddRelatedInfo(relatedInformation *Diagnostic) *Diagnostic {
	if relatedInformation != nil {
		d.RelatedInformation_ = append(d.RelatedInformation_, relatedInformation)
	}
	return d
}

func NewDiagnostic(file *SourceFile, loc core.TextRange, message *diagnostics.Message, args ...any) *Diagnostic {
	text := message.Message()
	if len(args) != 0 {
		text = core.FormatStringFromArgs(text, args)
	}
	return &Diagnostic{
		File_:     file,
		Loc_:      loc,
		Code_:     message.Code(),
		Category_: message.Category(),
		Message_:  text,
	}
}

// MessageChain

type MessageChain struct {
	Code_         int32
	Category_     diagnostics.Category
	Message_      string
	MessageChain_ []*MessageChain
}

func (m *MessageChain) Code() int32                    { return m.Code_ }
func (m *MessageChain) Category() diagnostics.Category { return m.Category_ }
func (m *MessageChain) Message() string                { return m.Message_ }
func (m *MessageChain) MessageChain() []*MessageChain  { return m.MessageChain_ }

func (m *MessageChain) AddMessageChain(messageChain *MessageChain) *MessageChain {
	if messageChain != nil {
		m.MessageChain_ = append(m.MessageChain_, messageChain)
	}
	return m
}

func NewMessageChain(message *diagnostics.Message, args ...any) *MessageChain {
	text := message.Message()
	if len(args) != 0 {
		text = core.FormatStringFromArgs(text, args)
	}
	return &MessageChain{
		Code_:     message.Code(),
		Category_: message.Category(),
		Message_:  text,
	}
}
