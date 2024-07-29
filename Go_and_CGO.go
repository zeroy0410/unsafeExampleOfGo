package main

/*
#include <stdio.h>

int cVar = 2333;
long long accessGoVariable(int* ptrGo) {
    int *ptrC = &cVar;
    // 在 C 代码中访问 Go 语言变量的地址
    printf("cVar addr: %p\n", ptrC);
    printf("goVar addr: %p\n", ptrGo);
    printf("Go variable value: %d\n", *ptrGo);
    // 尝试越界访问
    *(ptrGo) = 10;
    long long diff = ptrC - ptrGo;
    return diff; // 返回指针地址差值
}
void printCVar() {
    printf("C variable value after Go access: %d (from C function)", cVar);
}
*/
import "C"

import (
    "fmt"
    "unsafe"
)
func main() {
    var goVar int = 5
    fmt.Println("Go variable value before C access:", goVar)

    // 将 Go 变量的地址传递给 C 函数
    ptr := (*C.int)(unsafe.Pointer(&goVar))
    diff := C.accessGoVariable(ptr)
    fmt.Println("Difference between the addresses of a go variable and a C variable:", diff)

    // 输出被C语言修改后的 Go 变量的值
    fmt.Println("Go variable value after C access:", goVar)

    // Go 访问并修改C语言变量的值
    ptrC := (*int)(unsafe.Add(unsafe.Pointer(&goVar), diff * 4))
    fmt.Println("C variable value before Go access: ", *ptrC)
    *(ptrC) = 4666
    fmt.Println("C variable value after Go access:", *ptrC)
    C.printCVar()
}
// Go与C地址空间不隔离，C语言外部库可能破坏Go语言安全屏障，Go语言在unsafe环境下也可能对C语言地址空间造成影响。
// 当C语言外部库被劫持的情况下，Go语言的安全屏障将不起作用，敏感信息可被窃取。
