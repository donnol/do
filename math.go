package do

import (
	"math"
	"math/big"
)

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

// Factorial n <= 20, it will return 0 if n > 20, use FactorialBig instead.
func Factorial(n int) int {
	if n > 20 {
		return 0
	}

	s := 1
	for i := 2; i <= n; i++ {
		s *= i
	}
	return s
}

// FactorialBig n > 20
func FactorialBig(n int) string {
	s := big.NewInt(1)
	for i := 2; i <= n; i++ {
		s = s.Mul(s, big.NewInt(int64(i)))
	}
	return s.String()
}

// BinPow return a**b with binary pow. if `a` or `b` is very big, use `BinPowBig` instead
func BinPow(a, b int) int {
	res := 1
	for b > 0 {
		if b&1 != 0 {
			res = res * a
		}
		a = a * a
		b >>= 1
	}
	return res
}

// BinPow return a**b with binary pow
func BinPowBig(a, b int) string {
	res := big.NewInt(1)
	ai := big.NewInt(int64(a))
	for b > 0 {
		if b&1 != 0 {
			res = res.Mul(res, ai)
		}
		ai = ai.Mul(ai, ai)
		b >>= 1
	}
	return res.String()
}
