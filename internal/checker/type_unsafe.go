//go:build tsunsafe

package checker

import "unsafe"

type typeDataField struct {
	_dataType unsafe.Pointer
}

type interfaceValue struct {
	typ   unsafe.Pointer
	value unsafe.Pointer
}

func (t *typeDataField) data() typeData {
	var data typeData
	(*interfaceValue)(unsafe.Pointer(&data)).typ = t._dataType
	(*interfaceValue)(unsafe.Pointer(&data)).value = unsafe.Pointer(t)
	return data
}

func (t *typeDataField) setData(data typeData) {
	t._dataType = (*interfaceValue)(unsafe.Pointer(&data)).typ
}
