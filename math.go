package do

import "math"

func Round(x float64, p int) float64 {
	if p == 0 {
		return math.Round(x)
	}
	var n = 1.0
	for i := 0; i < p; i++ {
		n *= 10
	}
	return math.Round(x*n) / n
}

func Floor(x float64, p int) float64 {
	if p == 0 {
		return math.Floor(x)
	}
	var n = 1.0
	for i := 0; i < p; i++ {
		n *= 10
	}
	return math.Floor(x*n) / n
}

func Ceil(x float64, p int) float64 {
	if p == 0 {
		return math.Ceil(x)
	}
	var n = 1.0
	for i := 0; i < p; i++ {
		n *= 10
	}
	return math.Ceil(x*n) / n
}
