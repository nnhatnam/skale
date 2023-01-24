package trie

type Node[T Elem] struct {
	children map[T]*Node[T]
	lastElem bool
}

type Trie[T Elem] struct {
	root *Node[T]
	size int
}

func New[T Elem]() *Trie[T] {
	return &Trie[T]{root: &Node[T]{children: make(map[T]*Node[T])}}
}

func (t *Trie[T]) lazyInit() {
	if t.root == nil {
		t.root = &Node[T]{children: make(map[T]*Node[T])}
	}
}

func (t *Trie[T]) insert(s []T) {
	node := t.root
	for _, c := range s {
		if _, ok := node.children[c]; !ok {
			node.children[c] = &Node[T]{children: make(map[T]*Node[T])}
		}
		node = node.children[c]
	}

	if !node.lastElem {
		node.lastElem = true
		t.size++
	}
}

func (t *Trie[T]) findPrefix(s []T) bool {
	return findPrefix(t.root, s)
}

func findPrefix[T Elem](node *Node[T], s []T) bool {
	if node == nil {
		return false
	}

	if len(s) == 0 {
		return true
	}

	if _, ok := node.children[s[0]]; !ok {
		return false
	}

	return findPrefix(node.children[s[0]], s[1:])
}

//func (t *Trie[T]) get(elems []T) *Node[T] {
//	node := t.root
//	for _, elem := range elems {
//		if _, ok := node.children[elem]; !ok {
//			return nil
//		}
//		node = node.children[elem]
//	}
//	return node
//}

func (t *Trie[T]) get(elems []T) bool {
	node := t.root
	for _, elem := range elems {
		if _, ok := node.children[elem]; !ok {
			return false
		}
		node = node.children[elem]
	}

	return node.lastElem
}

func remove[T Elem](node *Node[T], elems []T) (_ *Node[T], deleted bool) {

	if node == nil {
		return node, false
	}

	if len(elems) == 0 && node.lastElem {

		if len(node.children) == 0 {
			return nil, true
		}
		node.lastElem = false
		return node, true
	}

	if len(elems) == 0 && !node.lastElem {
		return node, false
	}

	child, found := node.children[elems[0]]

	if !found {
		return node, false
	}

	node.children[elems[0]], deleted = remove[T](child, elems[1:])

	if node.children[elems[0]] == nil { // child is deleted
		delete(node.children, elems[0])

		if !node.lastElem && len(node.children) == 0 {
			return nil, true
		}

	}

	return node, deleted
}

func findAllPrefixFrom[T Elem](node *Node[T], elems []T, queue chan<- []T) {
	if node == nil {
		return
	}

	if node.lastElem {
		queue <- elems
	}

	for elem, child := range node.children {
		findAllPrefixFrom[T](child, append(elems, elem), queue)
	}

}

func (t *Trie[T]) getAllElems() (elems [][]T) {
	queue := make(chan []T)

	go func() {
		defer close(queue)
		findAllPrefixFrom(t.root, []T{}, queue)
	}()

	for elem := range queue {
		elems = append(elems, elem)
	}

	return elems
}

func (t *Trie[T]) getAllElemsWithPrefix(pre []T) (elems [][]T) {

	node := t.root
	for _, e := range pre {
		if _, ok := node.children[e]; !ok {
			return
		}

		node = node.children[e]
	}

	queue := make(chan []T)

	go func() {
		defer close(queue)
		findAllPrefixFrom(node, pre, queue)
	}()

	for elem := range queue {
		elems = append(elems, elem)
	}

	return elems

}

func (t *Trie[T]) longestPrefix(elems []T) (prefix []T) {
	node := t.root
	for _, elem := range elems {
		if _, ok := node.children[elem]; !ok {
			break
		}
		prefix = append(prefix, elem)
		node = node.children[elem]
	}

	return prefix
}

//func longestCommonPrefix[T Elem](a, b []T) (prefix []T) {
//	for i := 0; i < len(a) && i < len(b); i++ {
//		if a[i] != b[i] {
//			break
//		}
//		prefix = append(prefix, a[i])
//	}
//	return prefix
//}

func (t *Trie[T]) ElemCount() int {
	return t.size
}

func (t *Trie[T]) Insert(elems []T) {

	if len(elems) == 0 {
		return
	}
	t.insert(elems)
}

func (t *Trie[T]) Get(elems []T) bool {
	if len(elems) == 0 {
		return false
	}
	return t.get(elems)
}

func (t *Trie[T]) Delete(elems []T) (deleted bool) {
	if len(elems) == 0 {
		return
	}
	if _, deleted = remove[T](t.root, elems); deleted {
		t.size--
	}
	return deleted
}

func (t *Trie[T]) GetAll() [][]T {
	t.lazyInit()
	return t.getAllElems()
}

func (t *Trie[T]) GetAllWithPrefix(pre []T) [][]T {
	t.lazyInit()
	return t.getAllElemsWithPrefix(pre)
}

func (t *Trie[T]) LongestPrefix(elems []T) []T {
	t.lazyInit()
	return t.longestPrefix(elems)
}
