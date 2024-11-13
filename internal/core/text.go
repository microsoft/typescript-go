package core

// TextPos

type TextPos int32

// TextRange

type TextRange struct {
	Pos_ TextPos
	End_ TextPos
}

func NewTextRange(pos int, end int) TextRange {
	return TextRange{Pos_: TextPos(pos), End_: TextPos(end)}
}

func (t TextRange) Pos() int {
	return int(t.Pos_)
}

func (t TextRange) End() int {
	return int(t.End_)
}

func (t TextRange) Len() int {
	return int(t.End_ - t.Pos_)
}

func (t TextRange) ContainsInclusive(pos int) bool {
	return pos >= int(t.Pos_) && pos <= int(t.End_)
}
