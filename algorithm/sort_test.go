package algorithm

import (
	"fmt"
	"testing"
)

func TestQuickSort(t *testing.T) {
	nums := []int{2, 3, 1, 4, 5}
	lessFunc := func(a, b int) bool {
		return a < b
	}
	QuickSort[int](nums, lessFunc)
	for i := 0; i < len(nums); i++ {
		fmt.Println(nums[i])
	}
	alphas := []string{"bbc", "ab", "abed", "abd"}
	QuickSort[string](alphas, func(a, b string) bool {
		m, n := len(a), len(b)
		for i := 0; i < m && i < n; i++ {
			if a[i] < b[i] {
				return true
			} else if a[i] > b[i] {
				return false
			}
		}
		return m < n
	})
	for i := 0; i < len(alphas); i++ {
		fmt.Println(alphas[i])
	}
}
