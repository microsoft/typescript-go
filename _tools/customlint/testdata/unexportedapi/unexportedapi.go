package unexportedapi

type Foo struct {
	Bar *oops
}

type oops struct {
	v int
}

type Okay struct {
	Sure  int
	Value ***Okay2
}

type Okay2 struct {
	VeryGood struct{}
}

func OkayFunc(v *Okay) *Okay2 {
	if v == nil {
		return nil
	}
	return **v.Value
}
