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

// ref:https://www.bilibili.com/video/BV1jKpTeaEDB/?spm_id_from=333.337.search-card.all.click&vd_source=266c7f8a8bf830875ba1ad6a4063e7ca
// HeapSort 堆排序
// 堆的一些特性：
// 1. 最后一个非叶子节点的索引为heapSize/2-1. 因此如果索引大于heapSize/2-1, 说明该节点是叶子节点.
// 2. 节点i的父节点是(i-1)/2, 左子节点是2i+1, 右子节点是2i+2
// 3. 下层操作：siftDown, 堆化的核心操作, 也是堆排序的核心操作.
func HeapSort[T comparable](slice []T, less LessFunc[T]) {
	// 1. 堆化. 先把数组调整成堆. 从最后一个非叶子节点开始,逐个进行siftDown操作
	heapSize := len(slice)
	for i := heapSize/2 - 1; i >= 0; i-- {
		siftDown[T](slice, heapSize, i, less)
	}
	// 2. 基于堆进行排序. 堆化后堆顶元素一定是max/min, 因此把堆顶元素放到最后面, 并重新堆化前heapSize-1个元素.
	// 如果构建的是大根堆, 则排序后是升序. 反之.
	for end := heapSize - 1; end > 0; end-- {
		// 将堆顶元素与最后一个元素进行交换
		slice[end], slice[0] = slice[0], slice[end]
		// 重新堆化
		siftDown[T](slice, end, 0, less)
	}
}

// 把节点i进行下层操作. 下层就是当前节点不满足大根堆/小根堆时(除根节点外, 所有节点都大于/小于左右子节点),
// 将节点i与左右子节点中的最大/最小子节点进行交换.
func siftDown[T comparable](slice []T, heapSize, i int, less LessFunc[T]) {
	// 直到叶子节点才停止下层操作
	for i <= heapSize/2-1 {
		largerIndex := 2*i + 1
		rightIndex := 2*i + 2
		// 此处是把larger元素往上提, 因此构建的是大根堆.
		if rightIndex < heapSize && !less(slice[rightIndex], slice[largerIndex]) {
			largerIndex = rightIndex
		}
		// swap slice[i] and slice[smallerIndex]
		slice[i], slice[largerIndex] = slice[largerIndex], slice[i]
		i = largerIndex
	}
}
