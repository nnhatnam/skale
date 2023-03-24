package cartersian

import "github.com/nnhatnam/skale"

type Node[X, Y skale.Ordered] struct {
	left, right *Node[X, Y]

	x X // key or index
	y Y // priority
}

func NewNode[X, Y skale.Ordered](x X, y Y) *Node[X, Y] {
	return &Node[X, Y]{x: x, y: y}
}

type CartesianTree[X, Y skale.Ordered] struct {
	root *Node[X, Y]
}

func New[X, Y skale.Ordered]() *CartesianTree[X, Y] {
	return &CartesianTree[X, Y]{}
}

func merge[X, Y skale.Ordered](l, r *Node[X, Y]) *Node[X, Y] {
	// https://habr.com/en/post/101818/
	if l == nil {
		return r
	}

	if r == nil {
		return l
	}

	if l.y > r.y {
		l.right = merge[X, Y](l.right, r)
		return l
	}

	r.left = merge[X, Y](l, r.left)
	return r

}

// split splits the cartesian tree t into two trees l and r such that all keys in l are less than x and all keys in r are greater than or equal to x.
func split[X, Y skale.Ordered](t *Node[X, Y], x X) (*Node[X, Y], *Node[X, Y]) {
	// https://habr.com/en/post/101818/
	// http://e-maxx.ru/algo/treap
	if t == nil {
		return nil, nil
	}

	if t.x < x {
		l, r := split[X, Y](t.right, x)
		t.right = l
		return t, r
	}

	l, r := split[X, Y](t.left, x)
	t.left = r
	return l, t
}
