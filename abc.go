package main

import (
	"fmt"
	"math"
)

func main() {
	fmt.Println(math.MaxInt)
	fmt.Println(1<<31 - 1)

	var a = make([]int, 1, 1<<31-1)
	a[0] = 1
	fmt.Println(a[0])

}
