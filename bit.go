package do

func IsPowerOf2(n uint64) bool {
	// From: https://stackoverflow.com/questions/600293/how-to-check-if-a-number-is-a-power-of-2
	return n != 0 && n&(n-1) == 0
}

func IsSignDiff(x, y int) bool {
	return (x ^ y) < 0
}

func SwapByBit(a, b int) (int, int) {
	a ^= b
	b ^= a
	a ^= b
	return a, b
}
