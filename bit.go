package do

func IsPowerOf2(n uint64) bool {
	// From: https://stackoverflow.com/questions/600293/how-to-check-if-a-number-is-a-power-of-2
	return n != 0 && n&(n-1) == 0
}
