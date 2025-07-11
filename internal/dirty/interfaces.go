package dirty

type cloneable struct{}

func (c *cloneable) Clone() *cloneable {
	return &cloneable{}
}

type Cloneable[T any] interface {
	Clone() T
}

type Value[T any] interface {
	Value() T
	Original() T
	Dirty() bool
	Change(apply func(T))
	ChangeIf(cond func(T) bool, apply func(T)) bool
	Delete()
}
