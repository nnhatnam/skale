package main

import "github.com/nnhatnam/skale/bitmap"

func main() {
	b := bitmap.NewBitMap(1000000)
	for i := 0; i < 1000; i++ {
		b.Set(uint(i))
	}
}
