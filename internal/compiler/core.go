package compiler

import (
	"reflect"
	"strings"
)

func EquateStringsCaseInsensitive(a, b string) bool {
	// !!!
	// return a == b || strings.ToUpper(a) == strings.ToUpper(b)
	return strings.EqualFold(a, b)
}

func EquateStringsCaseSensitive(a, b string) bool {
	return a == b
}

func IsString(value interface{}) bool {
	_, ok := value.(string)
	return ok
}

func StartsWith(str string, prefix string, ignoreCase *bool) bool { //try it if the value returned is correct or not
	if *ignoreCase {
		return EquateStringsCaseInsensitive(str[0:len(prefix)], prefix)
	}
	return strings.LastIndex(str, prefix) == 0
}

func EndsWith(str string, suffix string, ignoreCase *bool) bool {
	expectedPos := len(str) - len(suffix)
	var result bool
	if *ignoreCase {
		result = EquateStringsCaseInsensitive(str[expectedPos:], suffix)
	} else {
		result = strings.Index(str, suffix) == expectedPos
		//result = str.indexOf(suffix, expectedPos) == expectedPos;
	}
	return expectedPos >= 0 && result
}

func IsSlice(value interface{}) bool {
	return reflect.TypeOf(value).Kind() == reflect.Slice
}
