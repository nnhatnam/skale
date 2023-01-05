package skale

type Ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64 | ~string
}

// Original design was creating an Orderable interface:
//
//	type Orderable[T any] interface {
//		Less(a T) int
//	}
//
// but this design force the user to implement the Less method for each type
// that she wants to use with the skale package. This is not a good design
// LessFunc would be a better design. The comparison function should only belong to the skale package.
type LessFunc[T any] func(a, b T) bool

func Less[T Ordered]() LessFunc[T] {
	return func(a, b T) bool {
		return a < b
	}
}
