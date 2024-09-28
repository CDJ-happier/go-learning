package algorithm

// LessFunc 比较函数
type LessFunc[T comparable] func(a, b T) bool

// QuickSort 实现泛型的slice[T]的快速排序
func QuickSort[T comparable](slice []T, less LessFunc[T]) {
	quickSort[T](slice, 0, len(slice)-1, less)
}

func quickSort[T comparable](slice []T, start, end int, less LessFunc[T]) {
	// 1. 终止条件
	//if left >= right {
	//	return
	//}
	// 2. 递归调用
	// 3. 结果整合
	if start < end {
		left, right := start, end
		pivot := slice[left] // 选择left 为轴点元素
		for left < right {
			// 从右往左找到比pivot小的元素
			for left < right && less(pivot, slice[right]) {
				right--
			}
			// 将slice[right] 放到left的位置, 即小于pivot的放到左边.
			// 此时right位置处于空闲, 且应该放大于pivot的元素. 即执行后一个for循环
			slice[left] = slice[right]
			// 从左往右找到比pivot大的元素
			for left < right && less(slice[left], pivot) {
				left++
			}
			// 将slice[left] 放到right的位置, 即大于pivot的放到右边
			// 此时left位置处于空闲, 且应该放小于pivot的元素. 即执行第一个for循环
			slice[right] = slice[left]
		}
		// 将pivot 放到left位置
		slice[left] = pivot
		quickSort(slice, start, left-1, less)
		quickSort(slice, left+1, end, less)
	}
}
