package skiplist_test

import (
	"fmt"
	"github.com/nnhatnam/skale/list/skiplist"
)

func ExampleSkipList() {
	l := skiplist.NewOrdered[int](64, 0.5)
	for i := 0; i < 10; i++ {
		l.ReplaceOrInsert(i)
	}
	fmt.Println("len:       ", l.Len())
	v, ok := l.Get(3)
	fmt.Println("get3:      ", v, ok)
	v, ok = l.Get(100)
	fmt.Println("get100:    ", v, ok)
	v, ok = l.Delete(4)
	fmt.Println("del4:      ", v, ok)
	v, ok = l.Delete(100)
	fmt.Println("del100:    ", v, ok)
	v, ok = l.ReplaceOrInsert(5)
	fmt.Println("replace5:  ", v, ok)
	v, ok = l.ReplaceOrInsert(100)
	fmt.Println("replace100:", v, ok)
	v, ok = l.Min()
	fmt.Println("min:       ", v, ok)
	v, ok = l.DeleteMin()
	fmt.Println("delmin:    ", v, ok)
	v, ok = l.Max()
	fmt.Println("max:       ", v, ok)
	v, ok = l.DeleteMax()
	fmt.Println("delmax:    ", v, ok)
	fmt.Println("len:       ", l.Len())
	// Output:
	// len:        10
	// get3:       3 true
	// get100:     0 false
	// del4:       4 true
	// del100:     0 false
	// replace5:   5 true
	// replace100: 0 false
	// min:        0 true
	// delmin:     0 true
	// max:        100 true
	// delmax:     100 true
	// len:        8
}
