package do

// SplitUint 0 -> [0], [1234] -> [1, 2, 3, 4]
func SplitUint(n uint64) []uint64 {
	if n == 0 {
		return []uint64{n}
	}

	parts := []uint64{}
	for n > 0 {
		parts = append(parts, n%10)
		n = n / 10
	}

	for i, j := 0, len(parts)-1; i < j; i, j = i+1, j-1 {
		parts[i], parts[j] = parts[j], parts[i]
	}

	return parts
}

func JoinUint(parts []uint64) (n uint64) {
	l := len(parts)
	if l == 0 {
		return
	}
	if l == 1 && parts[0] == 0 {
		return
	}
	for i := 0; i < l; i++ {
		n += parts[i] * Pow10(l-i-1)
	}
	return
}

// Pow10 10^n, return 0 if n < 0
func Pow10(n int) (r uint64) {
	if n < 0 {
		return 0
	}
	r = 1
	for i := 0; i < n; i++ {
		r *= 10
	}
	return
}
