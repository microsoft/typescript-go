package main

import (
	"errors"
	"strconv"

	"github.com/microsoft/typescript-go/internal/core"
)

var errParse = errors.New("parse error")

type tristateFlag core.Tristate

const (
 True string = "true"
 False       = "false"
 Unset       = "unset"
)

func (f *tristateFlag) Set(s string) error {
	v, err := strconv.ParseBool(s)
	if err != nil {
		return errParse
	}
	if v {
		*f = (tristateFlag)(core.TSTrue)
	} else {
		*f = (tristateFlag)(core.TSFalse)
	}
	return nil
}

func (f *tristateFlag) String() string {
	switch core.Tristate(*f) {
	case core.TSTrue:
		return True
	case core.TSFalse:
		return False
	default:
		return Unset
	}
}

func (f *tristateFlag) Get() any {
	return core.Tristate(*f)
}

func (f *tristateFlag) IsBoolFlag() bool {
	return true
}
