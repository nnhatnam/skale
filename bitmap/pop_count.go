package bitmap

// https://graphics.stanford.edu/~seander/bithacks.html#CountBitsSetParallel
// https://www.chessprogramming.org/Population_Count
func popCount(x uint64) uint64 {
	//TODO: writing explanation
	x = x - ((x >> 1) & 0x5555555555555555)
	x = (x & 0x3333333333333333) + ((x >> 2) & 0x3333333333333333)
	x = (x + (x >> 4)) & 0x0f0f0f0f0f0f0f0f
	return (x * 0x0101010101010101) >> 56
}
