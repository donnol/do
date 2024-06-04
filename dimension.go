package do

func Rectangle[T any](m, n int, initial ...T) [][]T {
	il := len(initial)

	dp := make([][]T, m)
	for i := 0; i < m; i++ {

		dp[i] = make([]T, n)
		for j := 0; j < n; j++ {
			var t T
			if il > 0 && il < n {
				t = initial[0]
			} else if il == n {
				t = initial[j]
			}
			dp[i][j] = t
		}

	}
	return dp
}

func Square[T any](n int, initial ...T) [][]T {
	return Rectangle[T](n, n, initial...)
}
