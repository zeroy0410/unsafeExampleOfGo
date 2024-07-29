package main

import (
	"fmt"
	"unsafe"
)

func main() {
	b := 2333
	a := [16]int{3: 3, 9: 9, 11: 11}
	// 获取变量a和b的地址
	addrA := unsafe.Pointer(&a[0])
	addrB := unsafe.Pointer(&b)

	diff := uintptr(addrB) - uintptr(addrA)

	fmt.Printf("Address of a: %p\n", addrA)
    fmt.Printf("Address of b: %p\n", addrB)
    fmt.Printf("Address difference: %d\n", diff)

	// unsafe包 越界访问
	ptr := (*int)(unsafe.Add(addrA, diff))
	fmt.Println(*ptr)
}