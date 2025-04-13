package main

import (
    "errors"
    "strconv"

    "github.com/microsoft/typescript-go/internal/core"
)

var errParse = errors.New("parse error")

type tristateFlag core.Tristate

func (f *tristateFlag) Set(s string) error {
    v, err := strconv.ParseBool(s)
    if err != nil {
        return errParse
    }
    *f = tristateFlag(core.TSTrue)
    if !v {
        *f = tristateFlag(core.TSFalse)
    }
    return nil
}

func (f *tristateFlag) String() string {
    return map[core.Tristate]string{
        core.TSTrue:  "true",
        core.TSFalse: "false",
    }[core.Tristate(*f)]
}

func (f *tristateFlag) Get() any {
    return *f
}

func (f *tristateFlag) IsBoolFlag() bool {
    return true
}


//
//
//
//    Changes made:
//
// * Simplifying logic in Set: Instead of using an if conditional structure, the default value is assigned and only modified in a specific case.
//
// * Using a Map in String: Replaced the switch with a Map to make returning values ​​based on states more concise and readable.
//
// * General cleanup: Redundancies were removed and code clarity was improved.
//
//
//
