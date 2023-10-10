package do

import "testing"

func TestRound(t *testing.T) {
	for _, cas := range []struct {
		In  float64
		P   []int
		Out []float64
	}{
		{2.4231313, []int{0, 1, 2}, []float64{2, 2.4, 2.42}},
		{2.5234242, []int{0, 1, 2}, []float64{3, 2.5, 2.52}},
		{3.4231313, []int{0, 1, 2}, []float64{3, 3.4, 3.42}},
		{3.5234242, []int{0, 1, 2}, []float64{4, 3.5, 3.52}},
	} {
		for i, p := range cas.P {
			r := Round(cas.In, p)
			Assert(t, r, cas.Out[i])
		}
	}
}

func TestFloor(t *testing.T) {
	for _, cas := range []struct {
		In  float64
		P   []int
		Out []float64
	}{
		{2.4231313, []int{0, 1, 2}, []float64{2, 2.4, 2.42}},
		{2.5234242, []int{0, 1, 2}, []float64{2, 2.5, 2.52}},
		{3.4231313, []int{0, 1, 2}, []float64{3, 3.4, 3.42}},
		{3.5234242, []int{0, 1, 2}, []float64{3, 3.5, 3.52}},
	} {
		for i, p := range cas.P {
			r := Floor(cas.In, p)
			Assert(t, r, cas.Out[i])
		}
	}
}

func TestCeil(t *testing.T) {
	for _, cas := range []struct {
		In  float64
		P   []int
		Out []float64
	}{
		{2.4231313, []int{0, 1, 2}, []float64{3, 2.5, 2.43}},
		{2.5234242, []int{0, 1, 2}, []float64{3, 2.6, 2.53}},
		{3.4231313, []int{0, 1, 2}, []float64{4, 3.5, 3.43}},
		{3.5234242, []int{0, 1, 2}, []float64{4, 3.6, 3.53}},
	} {
		for i, p := range cas.P {
			r := Ceil(cas.In, p)
			Assert(t, r, cas.Out[i])
		}
	}
}
