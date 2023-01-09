package avl

import (
	"fmt"
	"testing"
)

func TestAVLTree(t *testing.T) {
	avltree := NewOrdered[int]()
	avltree.Insert(5)
	avltree.Insert(8)
	avltree.Insert(9)
	avltree.Insert(9)
	avltree.Insert(9)
	avltree.Insert(9)
	avltree.Insert(9)
	avltree.Insert(9)

	str := "AVLTree\n"

	output[int](avltree.root, "", true, &str)
	fmt.Println(str)

}
