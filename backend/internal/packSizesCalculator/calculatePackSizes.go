package packSizesCalculator

import (
	"math"
	"slices"
)

func CalculatePackSizes(orderSize int, packSizes []int) map[int]int {
	if orderSize <= 0 || len(packSizes) == 0 {
		return nil
	}

	sorted := slices.Clone(packSizes)
	slices.SortFunc(sorted, func(a, b int) int { return b - a })

	largest := sorted[0]
	maxTarget := orderSize + largest

	dp := make([]int, maxTarget+1)
	for i := range dp {
		dp[i] = math.MaxInt32
	}
	dp[0] = 0

	for i := 1; i <= maxTarget; i++ {
		for _, pack := range sorted {
			if pack <= i && dp[i-pack] < math.MaxInt32 {
				if candidate := dp[i-pack] + 1; candidate < dp[i] {
					dp[i] = candidate
				}
			}
		}
	}

	bestTarget := orderSize
	for bestTarget <= maxTarget {
		if dp[bestTarget] < math.MaxInt32 {
			break
		}
		bestTarget++
	}

	result := make(map[int]int)
	remaining := bestTarget
	for remaining > 0 {
		for _, pack := range sorted {
			if pack <= remaining && dp[remaining-pack]+1 == dp[remaining] {
				result[pack]++
				remaining -= pack
				break
			}
		}
	}

	return result
}
