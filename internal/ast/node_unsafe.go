//go:build tsunsafe

package ast

import "unsafe"

type nodeDataField struct {
	_dataType unsafe.Pointer
}

type interfaceValue struct {
	typ   unsafe.Pointer
	value unsafe.Pointer
}

func (n *nodeDataField) data() nodeData {
	var data nodeData
	(*interfaceValue)(unsafe.Pointer(&data)).typ = n._dataType
	(*interfaceValue)(unsafe.Pointer(&data)).value = unsafe.Pointer(n)
	return data
}

func (n *nodeDataField) setData(data nodeData) {
	n._dataType = (*interfaceValue)(unsafe.Pointer(&data)).typ
}
