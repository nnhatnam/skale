package avl

import (
	"fmt"
	"testing"
)

func TestAVLTree(t *testing.T) {
	avltree := NewOrdered[int]()
	avltree.InsertNoReplace(5)
	avltree.InsertNoReplace(4)
	avltree.InsertNoReplace(7)
	avltree.InsertNoReplace(6)
	avltree.InsertNoReplace(8)
	str := "AVLTree\n"

	output[int](avltree.root, "", true, &str)
	fmt.Println(str)

}
