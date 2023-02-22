package bitmap

import "math/bits"

type BitMap struct {
	len    uint
	blocks []uint64
}

// NewBitMap returns a new BitMap with the given length.
func NewBitMap(len uint) *BitMap {

	return &BitMap{
		len:    len,
		blocks: make([]uint64, (len+63)>>6), // equivalent to (len + 64 - 1) / 64
	}
}

// Set sets a bit at the given index to 1.
func (b *BitMap) Set(i uint) *BitMap {
	if i >= b.len {
		panic("index out of range") // TODO: grow the bitmap if needed
	}
	// Concept through example:
	// Let supposed we have an array of 8 bit
	//   [0 0 0 0 0 0 0 0]
	// and we want to set the INDEX number 3 to 1 from right to left
	//   [0 0 0 0 1 0 0 0]  <- we want this
	// To achieve this, we shift 1 to the left by 3 bits, and then OR it with the original value
	//   [0 0 0 0 1 0 0 0]           <- 1 << 3 = 00001000
	// 			OR
	//   [0 0 0 0 0 0 0 0]           <- original value
	//   -----------------
	//   [0 0 0 0 1 0 0 0]           <- result
	blockIndex := i >> 6 // find block index, equivalent to i / 64
	bitIndex := i & 63   // find bit index in the block, equivalent to i % 64

	b.blocks[blockIndex] |= 1 << bitIndex // set the bit to 1
	return b
}

// Clear sets a bit at the given index to 0.
func (b *BitMap) Clear(i uint) *BitMap {
	if i >= b.len {
		panic("index out of range")
	}
	// Concept through example:
	// Let supposed we have an array of 8 bit
	//   [0 0 0 0 1 0 0 0]
	// and we want to set the INDEX number 3 to 0 from right to left
	//   [0 0 0 0 0 0 0 0]  <- we want this
	// To achieve this, we shift 1 to the left by 3 bits, and then inverts it (NOT), and then AND it with the original value
	//   [0 0 0 0 1 0 0 0]           <- original value
	// 			AND
	//   [1 1 1 1 0 1 1 1]           <- 1 << 3 = 00001000, then NOT = 11110111
	//   -----------------
	//   [0 0 0 0 0 0 0 0]           <- result
	blockIndex := i >> 6 // find block index, equivalent to i / 64
	bitIndex := i & 63   // find bit index in the block, equivalent to i % 64

	b.blocks[blockIndex] &^= 1 << bitIndex // set the bit to 0
	return b
}

func (b *BitMap) ClearRange(start, end uint) *BitMap {
	if start >= b.len || end >= b.len {
		panic("index out of range")
	}
	if start > end {
		panic("start index must be less than end index")
	}

	blockStart := start >> 6
	bitStart := start & 63

	blockEnd := end >> 6
	bitEnd := end & 63

	if blockStart == blockEnd {
		b.blocks[blockStart] &^= ((1 << (bitEnd + 1)) - 1) << bitStart
		return b
	}

	b.blocks[blockStart] &^= (1 << (64 - bitStart)) - 1
	for i := blockStart + 1; i < blockEnd; i++ {
		b.blocks[i] = 0
	}
	b.blocks[blockEnd] &^= (1 << (bitEnd + 1)) - 1

	return b
}

// ClearAll sets all bits to 0. Panics if the bitmap is empty.
func (b *BitMap) ClearAll() *BitMap {
	for i := range b.blocks {
		b.blocks[i] = 0
	}
	return b
}

// Get returns the value of a bit at the given index.
func (b *BitMap) Get(i uint) bool {
	if i >= b.len {
		panic("index out of range")
	}
	// To check if a particular bit index is set to 1 or not, we can use the AND operator
	// Concept through example:
	// Let supposed we have an array of 8 bit
	//   [1 0 0 0 1 0 0 0] 		<- original value, 1 << 3 = 00001000,
	// and we want to check if the INDEX number 3 is set to 1 or not from right to left
	// First, we shift 1 to the left by 3 bits, and then AND it with the original value.
	// If the result is not 0, then the bit is set to 1
	//   [1 0 0 0 1 0 0 0]           <- original value
	// 			AND
	//   [0 0 0 0 1 0 0 0]           <- 1 << 3 = 00001000
	//   -----------------
	//   [0 0 0 0 1 0 0 0]           <- result = 8 (00001000) != 0, so the bit is set to 1

	blockIndex := i >> 6 // find block index, equivalent to i / 64
	bitIndex := i & 63   // find bit index in the block, equivalent to i % 64

	return b.blocks[blockIndex]&(1<<bitIndex) != 0
}

// Flip flips the value of a bit at the given index.
func (b *BitMap) Flip(i uint) *BitMap {
	if i >= b.len {
		panic("index out of range")
	}
	// Concept through example:
	// Let supposed we have an array of 8 bit
	//   [0 0 0 0 1 0 0 0]
	// and we want to flip the INDEX number 3 from right to left
	//   [0 0 0 0 0 0 0 0]  <- we want this
	// To achieve this, we shift 1 to the left by 3 bits, and then XOR it with the original value
	//   [0 0 0 0 1 0 0 0]           <- original value
	// 			XOR
	//   [0 0 0 0 1 0 0 0]           <- 1 << 3 = 00001000
	//   -----------------
	//   [0 0 0 0 0 0 0 0]           <- result
	blockIndex := i >> 6 // find block index, equivalent to i / 64
	bitIndex := i & 63   // find bit index in the block, equivalent to i % 64

	b.blocks[blockIndex] ^= 1 << bitIndex // flip the bit
	return b
}

// IsEmpty returns true contains no bits sets to 1.
func (b *BitMap) IsEmpty() bool {
	for _, block := range b.blocks {
		if block != 0 {
			return false
		}
	}
	return true
}

// All returns true if all bits are set to 1. Panics if the bitmap is nil.
func (b *BitMap) All() bool {
	for _, block := range b.blocks {
		if block != ^uint64(0) { // TODO: investigate why this faster than using constant
			return false
		}
	}
	return true
}

// Any returns true if any bit is set to 1. Panics if the bitmap is nil.
func (b *BitMap) Any() bool {
	for _, block := range b.blocks {
		if block != 0 {
			return true
		}
	}
	return false
}

// PopCount returns the number of bits set to 1.
//func (b *BitMap) PopCount() uint64 {
//	var count uint64
//	for _, block := range b.blocks {
//		count += popCount(block)
//	}
//	return count
//}

// Count returns the number of bits set to 1.
func (b *BitMap) Count() uint64 {
	var count int
	for _, block := range b.blocks {
		count += bits.OnesCount64(block)
	}
	return uint64(count)
}

// And performs a bitwise AND operation between two bitmaps. Store the  result in the receiver b.
// Panics if the bitmaps have different lengths.
func (b *BitMap) And(other *BitMap) *BitMap {
	if b.len != other.len {
		panic("length mismatch")
	}
	for i := range b.blocks {
		b.blocks[i] &= other.blocks[i]
	}
	return b
}

// AndNot performs a bitwise AND NOT operation between two bitmaps. Store the  result in the receiver b.
// Panics if the bitmaps have different lengths.
func (b *BitMap) AndNot(other *BitMap) *BitMap {
	if b.len != other.len {
		panic("length mismatch")
	}
	for i := range b.blocks {
		b.blocks[i] &^= other.blocks[i]
	}
	return b
}

// Or performs a bitwise OR operation between two bitmaps. Store the  result in the receiver b.
// Panics if the bitmaps have different lengths.
func (b *BitMap) Or(other *BitMap) *BitMap {
	if b.len != other.len {
		panic("length mismatch")
	}
	for i := range b.blocks {
		b.blocks[i] |= other.blocks[i]
	}
	return b
}

// Xor performs a bitwise XOR operation between two bitmaps. Store the  result in the receiver b.
// Panics if the bitmaps have different lengths.
func (b *BitMap) Xor(other *BitMap) *BitMap {
	if b.len != other.len {
		panic("length mismatch")
	}
	for i := range b.blocks {
		b.blocks[i] ^= other.blocks[i]
	}
	return b
}

// Not performs a bitwise NOT operation on the bitmap. Store the  result in the receiver b.
func (b *BitMap) Not() *BitMap {
	for i := range b.blocks {
		b.blocks[i] = ^b.blocks[i]
	}
	return b
}

// Equal returns true if two bitmaps are equal.
func (b *BitMap) Equal(other *BitMap) bool {
	if b.len != other.len {
		return false
	}
	for i := range b.blocks {
		if b.blocks[i] != other.blocks[i] {
			return false
		}
	}
	return true
}
