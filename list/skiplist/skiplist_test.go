package skiplist

import (
	"golang.org/x/exp/slices"
	"math/rand"
	"reflect"
	"sort"
	"testing"
)

func TestNew(t *testing.T) {
	l := NewOrdered[int](64, 0.5)

	if l == nil {
		t.Error("New() failed")
	}

	if len(l.root.next) != 64 {
		t.Error("New() failed")
	}

	if l.root.next[0] == nil {
		t.Error("New() failed")
	}

	if l.root.next[0] != &l.root {
		t.Error("New() failed")
	}

	if l.root.prev == nil {
		t.Error("New() failed")
	}

	if l.root.prev != &l.root {
		t.Error("New() failed")
	}

	if l.maxLevel != 63 {
		t.Error("New() failed")
	}

	if l.p != 0.5 {
		t.Error("New() failed")
	}

	if l.less == nil {
		t.Error("New() failed")
	}

	if l.modifiedFingers == nil {
		t.Error("New() failed")
	}

	if len(l.modifiedFingers) != 64 {
		t.Error("New() failed")
	}

	if l.searchFingers == nil {
		t.Error("New() failed")
	}

	if len(l.searchFingers) != 64 {
		t.Error("New() failed")
	}

	if l.len != 0 {
		t.Error("New() failed")
	}

}

//func TestMaxLevel(t *testing.T) {
//	l := NewOrdered[int](64, 0.5)
//	var m = make(map[uint8]int)
//	for i := 0; i < 100000; i++ {
//		m[l.generateLevel()]++
//	}
//	fmt.Println(m)
//}

func getListItems[T any](l *SkipList[T]) []T {
	var items []T
	l.Ascend(func(item T) bool {
		items = append(items, item)
		return true
	})
	return items
}

func TestReplaceOrInsert(t *testing.T) {

	l := NewOrdered[int](64, 0.5)

	l.ReplaceOrInsert(1)

	if l.len != 1 {
		t.Errorf("ReplaceOrInsert failed: expected size 1, got %v", l.len)
	}

	if n, found := l.Get(1); !found {
		t.Errorf("ReplaceOrInsert failed: expected node with key 1, got %v and %v", n, found)
	}

	l.ReplaceOrInsert(2)

	if l.len != 2 {
		t.Errorf("ReplaceOrInsert failed: expected size 2, got %v", l.len)
	}

	if n, found := l.Get(2); !found {
		t.Errorf("ReplaceOrInsert failed: expected node with key 2, got %v and %v", n, found)
	}

	l.ReplaceOrInsert(3)
	l.ReplaceOrInsert(4)
	l.ReplaceOrInsert(5)

	items := getListItems(l)

	if !slices.Equal(items, []int{1, 2, 3, 4, 5}) {
		t.Errorf("ReplaceOrInsert failed: expected [1, 2, 3, 4, 5], got %v", items)
	}

	l.ReplaceOrInsert(10)
	l.ReplaceOrInsert(9)
	l.ReplaceOrInsert(8)
	l.ReplaceOrInsert(7)
	l.ReplaceOrInsert(6)

	items = getListItems(l)

	if !slices.Equal(items, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}) {
		t.Errorf("ReplaceOrInsert failed: expected [1, 2, 3, 4, 5, 6, 7, 8, 9, 10], got %v", items)
	}

	l.ReplaceOrInsert(1)

	items = getListItems(l)

	if !slices.Equal(items, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}) {
		t.Errorf("ReplaceOrInsert failed: expected [1, 2, 3, 4, 5, 6, 7, 8, 9, 10], got %v", items)
	}

	l2 := NewOrdered[int](63, 0.5)
	l2.ReplaceOrInsert(1)
	l2.ReplaceOrInsert(1)

	if l2.Len() != 1 {
		t.Errorf("expecting len 1")
	}
	if !l2.Has(1) {
		t.Errorf("expecting to find key=1")
	}

	l2.Delete(1)
	if l2.Len() != 0 {
		t.Errorf("expecting len 0")
	}
	if l2.Has(1) {
		t.Errorf("not expecting to find key=1")
	}

	l2.Delete(1)
	if l2.Len() != 0 {
		t.Errorf("expecting len 0")
	}

	if l2.Has(1) {
		t.Errorf("not expecting to find key=1")
	}

}

func TestReverseReplaceOrInsert(t *testing.T) {
	l := NewOrdered[int](64, 0.5)
	n := 100
	for i := 0; i < n; i++ {
		l.ReplaceOrInsert(n - i)
	}
	i := 0
	l.AscendGreaterOrEqual(0, func(item int) bool {
		i++
		if item != i {
			t.Errorf("bad order: got %d, expect %d", item, i)
		}
		return true
	})
}

func TestRange(t *testing.T) {
	l := NewOrdered[string](63, 0.5)
	order := []string{
		"ab", "aba", "abc", "a", "aa", "aaa", "b", "a-", "a!",
	}
	for _, i := range order {
		l.ReplaceOrInsert(i)
	}
	k := 0
	l.AscendRange("ab", "ac", func(item string) bool {
		if k > 3 {
			t.Fatalf("returned more items than expected")
		}
		i1 := order[k]
		i2 := item
		if i1 != i2 {
			t.Errorf("expecting %s, got %s", i1, i2)
		}
		k++
		return true
	})
}

func TestRandomInsertOrder(t *testing.T) {
	l := NewOrdered[int](16, 0.5)
	n := 1000
	perm := rand.Perm(n)
	for i := 0; i < n; i++ {
		l.ReplaceOrInsert(perm[i])
	}
	j := 0
	l.AscendGreaterOrEqual(0, func(item int) bool {

		if item != j {
			t.Fatalf("bad order %d, %d", item, j)
		}
		j++
		return true
	})
}

func TestRandomReplace(t *testing.T) {
	l := NewOrdered[int](48, 0.5)
	n := 100
	perm := rand.Perm(n)
	for i := 0; i < n; i++ {
		l.ReplaceOrInsert(perm[i])
	}
	perm = rand.Perm(n)
	for i := 0; i < n; i++ {
		if replaced, isReplace := l.ReplaceOrInsert(perm[i]); !isReplace || replaced != perm[i] {
			t.Errorf("error replacing")
		}
	}
}

func TestRandomInsertSequentialDelete(t *testing.T) {
	l := NewOrdered[int](16, 0.5)

	n := 1000
	perm := rand.Perm(n)
	for i := 0; i < n; i++ {
		l.ReplaceOrInsert(perm[i])
	}

	for i := 0; i < n; i++ {
		l.Delete(i)
	}

}

func TestRandomInsertDeleteNonExistent(t *testing.T) {
	l := NewOrdered[int](8, 0.5)
	n := 100
	perm := rand.Perm(n)
	for i := 0; i < n; i++ {
		l.ReplaceOrInsert(perm[i])
	}

	for i := 0; i < n; i++ {
		if u, success := l.Delete(i); !success || u != i {
			t.Errorf("delete failed")
		}
	}

	if _, success := l.Delete(200); success {
		t.Errorf("deleted non-existent item")
	}
	if _, success := l.Delete(-2); success {
		t.Errorf("deleted non-existent item")
	}

	if _, success := l.Delete(200); success {
		t.Errorf("deleted non-existent item")
	}
	if _, success := l.Delete(-2); success {
		t.Errorf("deleted non-existent item")
	}
}

func TestRandomInsertPartialDeleteOrder(t *testing.T) {
	l := NewOrdered[int](8, 0.5)
	n := 100
	perm := rand.Perm(n)
	for i := 0; i < n; i++ {
		l.ReplaceOrInsert(perm[i])
	}

	for i := 1; i < n-1; i++ {
		l.Delete(i)

	}

	j := 0
	l.AscendGreaterOrEqual(0, func(item int) bool {
		switch j {
		case 0:
			if item != 0 {
				t.Errorf("expecting 0")
			}
		case 1:
			if item != n-1 {
				t.Errorf("expecting %d", n-1)
			}
		}
		j++
		return true
	})
}

func TestInsertNoReplace(t *testing.T) {
	l := NewOrdered[int](8, 0.5)
	n := 1000
	for q := 0; q < 2; q++ {
		perm := rand.Perm(n)
		for i := 0; i < n; i++ {
			l.InsertNoReplace(perm[i])
		}
	}

	j := 0
	l.AscendGreaterOrEqual(0, func(item int) bool {
		if item != j/2 {
			t.Fatalf("bad order")
		}
		j++
		return true
	})
}

// Test case from https://github.com/google/btree

func intRange(s int, reverse bool) []int {
	out := make([]int, s)
	for i := 0; i < s; i++ {
		v := i
		if reverse {
			v = s - i - 1
		}
		out[i] = v
	}
	return out
}

func intAll(l *SkipList[int]) (out []int) {
	l.Ascend(func(a int) bool {
		out = append(out, a)
		return true
	})
	return
}

func intAllRev(l *SkipList[int]) (out []int) {
	l.Descend(func(a int) bool {
		out = append(out, a)
		return true
	})
	return
}

func TestSkipListG(t *testing.T) {
	l := NewOrdered[int](8, 0.5)
	const lSize = 10000
	//const lSize = 30
	for i := 0; i < 10; i++ {
		if min, ok := l.Min(); ok || min != 0 {
			t.Fatalf("empty min, got %+v", min)
		}

		if max, ok := l.Max(); ok || max != 0 {
			t.Fatalf("empty max, got %+v", max)
		}
		for _, item := range rand.Perm(lSize) {
			//fmt.Println("????? insert ????: ", item)
			if x, ok := l.ReplaceOrInsert(item); ok || x != 0 {
				t.Fatal("insert found item", item)
			}
		}

		//fmt.Println("-------------------------------------------------------------------------------")
		for _, item := range rand.Perm(lSize) {
			//fmt.Println("????? insert again ????: ", item)
			if x, ok := l.ReplaceOrInsert(item); !ok || x != item {
				t.Fatal("insert didn't find item", item)
			}
		}

		want := 0
		if min, ok := l.Min(); !ok || min != want {
			t.Fatalf("min: ok %v want %+v, got %+v", ok, want, min)
		}
		want = lSize - 1
		if max, ok := l.Max(); !ok || max != want {
			t.Fatalf("max: ok %v want %+v, got %+v", ok, want, max)
		}

		got := intAll(l)
		wantRange := intRange(lSize, false)
		if !reflect.DeepEqual(got, wantRange) {
			t.Fatalf("mismatch:\n got: %v\nwant: %v", got, wantRange)
		}

		gotrev := intAllRev(l)
		wantrev := intRange(lSize, true)
		if !reflect.DeepEqual(gotrev, wantrev) {
			t.Fatalf("mismatch:\n got: %v\nwant: %v", gotrev, wantrev)
		}

		l.NodeAscend(func(item *Node[int]) bool {

			for _, ptr := range item.next {
				if ptr == nil {
					t.Fatalf("there is nil pointer in next array")
				}
			}

			return true
		})

		for _, item := range rand.Perm(lSize) {
			//fmt.Println("delete: ", item, lSize)
			if x, ok := l.Delete(item); !ok || x != item {
				t.Fatalf("didn't find %v", item)
			}
		}

		if got = intAll(l); len(got) > 0 {
			t.Fatalf("some left!: %v", got)
		}
		if got = intAllRev(l); len(got) > 0 {
			t.Fatalf("some left!: %v", got)
		}

	}
}

func TestDeleteMinG(t *testing.T) {
	l := NewOrdered[int](64, 0.5)
	for _, v := range rand.Perm(100) {
		l.ReplaceOrInsert(v)
	}
	var got []int
	for v, ok := l.DeleteMin(); ok; v, ok = l.DeleteMin() {
		got = append(got, v)
	}
	if want := intRange(100, false); !reflect.DeepEqual(got, want) {
		t.Fatalf("ascendrange:\n got: %v\nwant: %v", got, want)
	}
}

func TestDeleteMaxG(t *testing.T) {
	l := NewOrdered[int](64, 0.5)
	for _, v := range rand.Perm(100) {
		l.ReplaceOrInsert(v)
	}
	var got []int
	for v, ok := l.DeleteMax(); ok; v, ok = l.DeleteMax() {
		got = append(got, v)
	}
	if want := intRange(100, true); !reflect.DeepEqual(got, want) {
		t.Fatalf("ascendrange:\n got: %v\nwant: %v", got, want)
	}
}

func TestAscendRangeG(t *testing.T) {
	l := NewOrdered[int](6, 0.5)
	for _, v := range rand.Perm(100) {
		l.ReplaceOrInsert(v)
	}
	var got []int
	l.AscendRange(40, 60, func(a int) bool {
		got = append(got, a)
		return true
	})
	if want := intRange(100, false)[40:60]; !reflect.DeepEqual(got, want) {
		t.Fatalf("ascendrange:\n got: %v\nwant: %v", got, want)
	}
	got = got[:0]
	l.AscendRange(40, 60, func(a int) bool {
		if a > 50 {
			return false
		}
		got = append(got, a)
		return true
	})
	if want := intRange(100, false)[40:51]; !reflect.DeepEqual(got, want) {
		t.Fatalf("ascendrange:\n got: %v\nwant: %v", got, want)
	}
}

func TestDescendRangeG(t *testing.T) {
	l := NewOrdered[int](6, 0.5)
	for _, v := range rand.Perm(100) {
		l.ReplaceOrInsert(v)
	}
	var got []int
	l.DescendRange(60, 40, func(a int) bool {
		got = append(got, a)
		return true
	})
	if want := intRange(100, true)[39:59]; !reflect.DeepEqual(got, want) {
		t.Fatalf("descendrange:\n got: %v\nwant: %v", got, want)
	}
	got = got[:0]
	l.DescendRange(60, 40, func(a int) bool {
		if a < 50 {
			return false
		}
		got = append(got, a)
		return true
	})
	if want := intRange(100, true)[39:50]; !reflect.DeepEqual(got, want) {
		t.Fatalf("descendrange:\n got: %v\nwant: %v", got, want)
	}
}

func TestAscendLessThanG(t *testing.T) {
	l := NewOrdered[int](6, 0.5)
	for _, v := range rand.Perm(100) {
		l.ReplaceOrInsert(v)
	}
	var got []int
	l.AscendLessThan(60, func(a int) bool {
		got = append(got, a)
		return true
	})
	if want := intRange(100, false)[:60]; !reflect.DeepEqual(got, want) {
		t.Fatalf("ascendrange:\n got: %v\nwant: %v", got, want)
	}
	got = got[:0]
	l.AscendLessThan(60, func(a int) bool {
		if a > 50 {
			return false
		}
		got = append(got, a)
		return true
	})
	if want := intRange(100, false)[:51]; !reflect.DeepEqual(got, want) {
		t.Fatalf("ascendrange:\n got: %v\nwant: %v", got, want)
	}
}

func TestDescendLessOrEqualG(t *testing.T) {
	l := NewOrdered[int](6, 0.5)
	for _, v := range rand.Perm(100) {
		l.ReplaceOrInsert(v)
	}
	var got []int
	l.DescendLessOrEqual(40, func(a int) bool {
		got = append(got, a)
		return true
	})
	if want := intRange(100, true)[59:]; !reflect.DeepEqual(got, want) {
		t.Fatalf("descendlessorequal:\n got: %v\nwant: %v", got, want)
	}
	got = got[:0]
	l.DescendLessOrEqual(60, func(a int) bool {
		if a < 50 {
			return false
		}
		got = append(got, a)
		return true
	})
	if want := intRange(100, true)[39:50]; !reflect.DeepEqual(got, want) {
		t.Fatalf("descendlessorequal:\n got: %v\nwant: %v", got, want)
	}
}

func TestAscendGreaterOrEqualG(t *testing.T) {
	l := NewOrdered[int](6, 0.5)
	for _, v := range rand.Perm(100) {
		l.ReplaceOrInsert(v)
	}
	var got []int
	l.AscendGreaterOrEqual(40, func(a int) bool {
		got = append(got, a)
		return true
	})
	if want := intRange(100, false)[40:]; !reflect.DeepEqual(got, want) {
		t.Fatalf("ascendrange:\n got: %v\nwant: %v", got, want)
	}
	got = got[:0]
	l.AscendGreaterOrEqual(40, func(a int) bool {
		if a > 50 {
			return false
		}
		got = append(got, a)
		return true
	})
	if want := intRange(100, false)[40:51]; !reflect.DeepEqual(got, want) {
		t.Fatalf("ascendrange:\n got: %v\nwant: %v", got, want)
	}
}

func TestDescendGreaterThanG(t *testing.T) {
	l := NewOrdered[int](6, 0.5)
	for _, v := range rand.Perm(100) {
		l.ReplaceOrInsert(v)
	}
	var got []int
	l.DescendGreaterThan(40, func(a int) bool {
		got = append(got, a)
		return true
	})
	if want := intRange(100, true)[:59]; !reflect.DeepEqual(got, want) {
		t.Fatalf("descendgreaterthan:\n got: %v\nwant: %v", got, want)
	}
	got = got[:0]
	l.DescendGreaterThan(40, func(a int) bool {
		if a < 50 {
			return false
		}
		got = append(got, a)
		return true
	})
	if want := intRange(100, true)[:50]; !reflect.DeepEqual(got, want) {
		t.Fatalf("descendgreaterthan:\n got: %v\nwant: %v", got, want)
	}
}

const benchmarkListSize = 1000000

func BenchmarkInsertG(b *testing.B) {
	b.StopTimer()
	insertP := rand.Perm(benchmarkListSize)
	b.StartTimer()
	i := 0
	for i < b.N {
		l := NewOrdered[int](64, 0.5)
		for _, item := range insertP {
			l.ReplaceOrInsert(item)
			i++
			if i >= b.N {
				return
			}
		}
	}
}

func BenchmarkSeekG(b *testing.B) {
	b.StopTimer()
	size := 100000
	insertP := rand.Perm(size)
	l := NewOrdered[int](64, 0.5)
	for _, item := range insertP {
		l.ReplaceOrInsert(item)
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		l.AscendGreaterOrEqual(i%size, func(i int) bool { return false })
	}
}

func BenchmarkDeleteInsertG(b *testing.B) {
	b.StopTimer()
	insertP := rand.Perm(benchmarkListSize)
	l := NewOrdered[int](64, 0.5)
	for _, item := range insertP {
		l.ReplaceOrInsert(item)
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		l.Delete(insertP[i%benchmarkListSize])
		l.ReplaceOrInsert(insertP[i%benchmarkListSize])
	}
}

//func BenchmarkDeleteInsertCloneOnceG(b *testing.B) {
//	b.StopTimer()
//	insertP := rand.Perm(benchmarkListSize)
//	l := NewOrdered[int](64, 0.5)
//	for _, item := range insertP {
//		l.ReplaceOrInsert(item)
//	}
//	tr = l.Clone()
//	b.StartTimer()
//	for i := 0; i < b.N; i++ {
//		l.Delete(insertP[i%benchmarkListSize])
//		l.ReplaceOrInsert(insertP[i%benchmarkListSize])
//	}
//}

//func BenchmarkDeleteInsertCloneEachTimeG(b *testing.B) {
//	b.StopTimer()
//	insertP := rand.Perm(benchmarkListSize)
//	l := NewOrdered[int](64, 0.5)
//	for _, item := range insertP {
//		l.ReplaceOrInsert(item)
//	}
//	b.StartTimer()
//	for i := 0; i < b.N; i++ {
//		tr = l.Clone()
//		l.Delete(insertP[i%benchmarkListSize])
//		l.ReplaceOrInsert(insertP[i%benchmarkListSize])
//	}
//}

func BenchmarkDeleteG(b *testing.B) {
	b.StopTimer()
	insertP := rand.Perm(benchmarkListSize)
	removeP := rand.Perm(benchmarkListSize)
	b.StartTimer()
	i := 0
	for i < b.N {
		b.StopTimer()
		l := NewOrdered[int](64, 0.5)
		for _, v := range insertP {
			l.ReplaceOrInsert(v)
		}
		b.StartTimer()
		for _, item := range removeP {
			l.Delete(item)
			i++
			if i >= b.N {
				return
			}
		}
		if l.Len() > 0 {
			panic(l.Len())
		}
	}
}

func BenchmarkGetG(b *testing.B) {
	b.StopTimer()
	insertP := rand.Perm(benchmarkListSize)
	removeP := rand.Perm(benchmarkListSize)
	b.StartTimer()
	i := 0
	for i < b.N {
		b.StopTimer()
		l := NewOrdered[int](64, 0.5)
		for _, v := range insertP {
			l.ReplaceOrInsert(v)
		}
		b.StartTimer()
		for _, item := range removeP {
			l.Get(item)
			i++
			if i >= b.N {
				return
			}
		}
	}
}

//func BenchmarkGetCloneEachTimeG(b *testing.B) {
//	b.StopTimer()
//	insertP := rand.Perm(benchmarkListSize)
//	removeP := rand.Perm(benchmarkListSize)
//	b.StartTimer()
//	i := 0
//	for i < b.N {
//		b.StopTimer()
//		l := NewOrdered[int](64, 0.5)
//		for _, v := range insertP {
//			l.ReplaceOrInsert(v)
//		}
//		b.StartTimer()
//		for _, item := range removeP {
//			tr = l.Clone()
//			l.Get(item)
//			i++
//			if i >= b.N {
//				return
//			}
//		}
//	}
//}

func BenchmarkAscendG(b *testing.B) {
	arr := rand.Perm(benchmarkListSize)
	l := NewOrdered[int](64, 0.5)
	for _, v := range arr {
		l.ReplaceOrInsert(v)
	}
	sort.Ints(arr)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		j := 0
		l.Ascend(func(item int) bool {
			if item != arr[j] {
				b.Fatalf("mismatch: expected: %v, got %v", arr[j], item)
			}
			j++
			return true
		})
	}
}

func BenchmarkDescendG(b *testing.B) {
	arr := rand.Perm(benchmarkListSize)
	l := NewOrdered[int](64, 0.5)
	for _, v := range arr {
		l.ReplaceOrInsert(v)
	}
	sort.Ints(arr)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		j := len(arr) - 1
		l.Descend(func(item int) bool {
			if item != arr[j] {
				b.Fatalf("mismatch: expected: %v, got %v", arr[j], item)
			}
			j--
			return true
		})
	}
}

func BenchmarkAscendRangeG(b *testing.B) {
	arr := rand.Perm(benchmarkListSize)
	l := NewOrdered[int](64, 0.5)
	for _, v := range arr {
		l.ReplaceOrInsert(v)
	}
	sort.Ints(arr)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		j := 100
		l.AscendRange(100, arr[len(arr)-100], func(item int) bool {
			if item != arr[j] {
				b.Fatalf("mismatch: expected: %v, got %v", arr[j], item)
			}
			j++
			return true
		})
		if j != len(arr)-100 {
			b.Fatalf("expected: %v, got %v", len(arr)-100, j)
		}
	}
}

func BenchmarkDescendRangeG(b *testing.B) {
	arr := rand.Perm(benchmarkListSize)
	l := NewOrdered[int](64, 0.5)
	for _, v := range arr {
		l.ReplaceOrInsert(v)
	}
	sort.Ints(arr)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		j := len(arr) - 100
		l.DescendRange(arr[len(arr)-100], 100, func(item int) bool {
			if item != arr[j] {
				b.Fatalf("mismatch: expected: %v, got %v", arr[j], item)
			}
			j--
			return true
		})
		if j != 100 {
			b.Fatalf("expected: %v, got %v", len(arr)-100, j)
		}
	}
}

func BenchmarkAscendGreaterOrEqualG(b *testing.B) {
	arr := rand.Perm(benchmarkListSize)
	l := NewOrdered[int](64, 0.5)
	for _, v := range arr {
		l.ReplaceOrInsert(v)
	}
	sort.Ints(arr)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		j := 100
		k := 0
		l.AscendGreaterOrEqual(100, func(item int) bool {
			if item != arr[j] {
				b.Fatalf("mismatch: expected: %v, got %v", arr[j], item)
			}
			j++
			k++
			return true
		})
		if j != len(arr) {
			b.Fatalf("expected: %v, got %v", len(arr), j)
		}
		if k != len(arr)-100 {
			b.Fatalf("expected: %v, got %v", len(arr)-100, k)
		}
	}
}

func BenchmarkDescendLessOrEqualG(b *testing.B) {
	arr := rand.Perm(benchmarkListSize)
	l := NewOrdered[int](64, 0.5)
	for _, v := range arr {
		l.ReplaceOrInsert(v)
	}
	sort.Ints(arr)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		j := len(arr) - 100
		k := len(arr)
		l.DescendLessOrEqual(arr[len(arr)-100], func(item int) bool {
			if item != arr[j] {
				b.Fatalf("mismatch: expected: %v, got %v", arr[j], item)
			}
			j--
			k--
			return true
		})
		if j != -1 {
			b.Fatalf("expected: %v, got %v", -1, j)
		}
		if k != 99 {
			b.Fatalf("expected: %v, got %v", 99, k)
		}
	}
}
