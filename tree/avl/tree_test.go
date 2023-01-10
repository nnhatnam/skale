package avl

import (
	"fmt"
	"testing"
)

func TestAVLTree(t *testing.T) {
	avltree := NewOrdered[int]()
	avltree.Insert(5)
	avltree.Insert(4)
	avltree.Insert(7)
	avltree.Insert(6)
	avltree.Insert(8)
	str := "AVLTree\n"

	output[int](avltree.root, "", true, &str)
	fmt.Println(str)

}
