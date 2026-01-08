package heap

import (
	"errors"
	"sort"
)

var errPositionOutOfRange = errors.New("position out of range")

func FindKthPreferred(ratings []int, position int) (int, error) {
	if position <= 0 || position > len(ratings) {
		return 0, errPositionOutOfRange
	}

	sorted := make([]int, len(ratings))
	copy(sorted, ratings)

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i] > sorted[j]
	})

	return sorted[position-1], nil
}
