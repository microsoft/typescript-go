//go:build !tsunsafe

package ast

type nodeDataField struct {
	_data nodeData
}

func (n *nodeDataField) data() nodeData {
	return n._data
}

func (n *nodeDataField) setData(data nodeData) {
	n._data = data
}
