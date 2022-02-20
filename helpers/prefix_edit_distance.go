package helpers

func PrefixEditDistance(str string, prefix string) [][]int {
	dp := [][]int{}

	for i := 0; i <= len(str); i++ {
		dp = append(dp, make([]int, len(prefix)+1))
	}

	for j := 0; j <= len(prefix); j++ {
		dp[0][j] = j
	}

	for i := 1; i <= len(str); i++ {
		for j := 1; j <= len(prefix); j++ {
			// result := 0
			// if i > j {
			// 	dp[i][j] = dp[i-1][j]
			// } else if i < j {
			// 	if9
			// }
			// if str[i-1] == prefix[j-1] {
			// 	dp[i][j] = dp[i-1][j-1]
			// } else {
			// 	if i == j {

			// 	}
			// 	min := dp[i-1][j-1]

			// 	if dp[i-1][j] < min {
			// 		min = dp[i-1][j]
			// 	}

			// 	if dp[i][j-1] < min {
			// 		min = dp[i][j-1]
			// 	}

			// 	dp[i][j] = min + 1
			// }
		}
	}

	return dp
}
