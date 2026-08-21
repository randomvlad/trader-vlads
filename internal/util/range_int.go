package util

type RangeInt struct {
	valueMin int
	valueMax int
}

func NewRangeInt(valueMin, valueMax int) *RangeInt {
	return &RangeInt{
		valueMin: valueMin,
		valueMax: valueMax,
	}
}

func (s *RangeInt) Generate(r *RandomGenerator) int {
	return r.RollInt(s.valueMin, s.valueMax)
}
