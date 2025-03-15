//go:build !tsunsafe

package checker

type typeDataField struct {
	_data typeData
}

func (t *typeDataField) data() typeData {
	return t._data
}

func (t *typeDataField) setData(data typeData) {
	t._data = data
}
