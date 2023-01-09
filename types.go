package skale

type Ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64 | ~string
}

// LessFunc[T] determines how to order a type 'T'
type LessFunc[T any] func(a, b T) bool

func Less[T Ordered]() LessFunc[T] {
	return func(a, b T) bool {
		return a < b
	}
}
