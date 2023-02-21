package bitmap

type BitMap struct {
	len    uint
	blocks []uint64
}

// NewBitMap returns a new BitMap with the given length.
func NewBitMap(len uint) *BitMap {
	return &BitMap{
		len:    len,
		blocks: make([]uint64, (len+63)>>6),
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
