package main

import "fmt"

func main() {
	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println(r)
				fmt.Println("recovered from panic in outer")
			}
		}()
		invokePanic()              // 内部如果没有使用recover处理panic, 那么对于当前匿名函数来说, invokePanic()是一个panic
		fmt.Println("after panic") // 不会执行
	}()
}

func invokePanic() {
	defer fmt.Println("after invoke panic")
	//defer func() {
	//	if r := recover(); r != nil {
	//		fmt.Println("recovered from panic in inner")
	//	}
	//}()
	panic("panic")
}
