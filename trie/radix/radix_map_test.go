package radix

import (
	crand "crypto/rand"
	"fmt"
	"golang.org/x/exp/slices"
	"reflect"
	"sort"
	"testing"
)

// borrow test cases from github.com/armon/go-radix
func generateUUID() string {
	buf := make([]byte, 16)
	if _, err := crand.Read(buf); err != nil {
		panic(fmt.Errorf("failed to read random bytes: %v", err))
	}

	return fmt.Sprintf("%08x-%04x-%04x-%04x-%12x",
		buf[0:4],
		buf[4:6],
		buf[6:8],
		buf[8:10],
		buf[10:16])
}

func TestNewRadixMap(t *testing.T) {
	var min, max string
	inp := make(map[string]int)
	for i := 0; i < 10000; i++ {
		gen := generateUUID()
		inp[gen] = i
		if gen < min || i == 0 {
			min = gen
		}
		if gen > max || i == 0 {
			max = gen
		}
	}

	r := NewRadixMap[byte, int]()

	for k, v := range inp {
		r.ReplaceOrInsert([]byte(k), v)
	}

	if r.Len() != len(inp) {
		t.Errorf("Got %v expected %v", r.Len(), len(inp))
	}

	r.AscendGreaterOrEqual([]byte(min), func(key []byte, value int) bool {
		if string(key) < min {
			t.Errorf("Got %v expected %v", string(key), min)
		}
		if string(key) > max {
			t.Errorf("Got %v expected %v", string(key), max)
		}
		return false
	})

	for k, v := range inp {
		out, ok := r.Get([]byte(k))

		if !ok {
			t.Fatalf("missing key: %v", k)
		}
		if out != v {
			t.Fatalf("value mis-match: %v %v", out, v)
		}
	}

	// Check min and max
	outMin, _, _ := r.Min()
	if string(outMin) != min {
		t.Fatalf("bad minimum: %v %v", outMin, min)
	}
	outMax, _, _ := r.Max()
	if string(outMax) != max {
		t.Fatalf("bad maximum: %v %v", outMax, max)
	}

	for k, v := range inp {
		//fmt.Println("deleting: ", k, v)
		out, ok := r.Delete([]byte(k))
		//fmt.Println("deleted: ", k, v, out, ok)
		if !ok {
			t.Fatalf("missing key: %v", k)
		}
		if out != v {
			t.Fatalf("value mis-match: %v %v", out, v)
		}
	}
	if r.Len() != 0 {
		t.Fatalf("bad length: %v", r.Len())
	}

}

func TestNewRadixMapExtend(t *testing.T) {
	var min, max string
	inp := make(map[string]int)
	for i := 0; i < 10000; i++ {
		gen := generateUUID()
		inp[gen] = i
		if gen < min || i == 0 {
			min = gen
		}
		if gen > max || i == 0 {
			max = gen
		}
	}

	r := NewRadixMap[byte, int]()

	for k, v := range inp {
		r.ReplaceOrInsertExtend([]byte(k), v, nil)
	}

	if r.Len() != len(inp) {
		t.Errorf("Got %v expected %v", r.Len(), len(inp))
	}

	r.AscendGreaterOrEqual([]byte(min), func(key []byte, value int) bool {
		if string(key) < min {
			t.Errorf("Got %v expected %v", string(key), min)
		}
		if string(key) > max {
			t.Errorf("Got %v expected %v", string(key), max)
		}
		return false
	})

	for k, v := range inp {
		out, ok := r.Get([]byte(k))

		if !ok {
			t.Fatalf("missing key: %v", k)
		}
		if out != v {
			t.Fatalf("value mis-match: %v %v", out, v)
		}

		r.replaceOrInsertExtend([]byte(k), v, func(oldValue int) int {
			return oldValue + v
		})

		out, ok = r.Get([]byte(k))

		if !ok {
			t.Fatalf("missing key: %v", k)
		}
		if out != v*2 {
			t.Fatalf("value mis-match: %v %v", out, v)
		}
	}

	// Check min and max
	outMin, _, _ := r.Min()
	if string(outMin) != min {
		t.Fatalf("bad minimum: %v %v", outMin, min)
	}
	outMax, _, _ := r.Max()
	if string(outMax) != max {
		t.Fatalf("bad maximum: %v %v", outMax, max)
	}

	for k, v := range inp {
		out, ok := r.Delete([]byte(k))
		if !ok {
			t.Fatalf("missing key: %v", k)
		}
		if out != v*2 {
			t.Fatalf("value mis-match: %v %v", out, v)
		}
	}
	if r.Len() != 0 {
		t.Fatalf("bad length: %v", r.Len())
	}

}

func TestEmptyKey(t *testing.T) {
	r := NewRadixMap[byte, bool]()

	s := []string{"", "A", "AB"}

	for _, ss := range s {
		r.ReplaceOrInsert([]byte(ss), true)
	}

	var inOrder []string

	r.AscendGreaterOrEqual([]byte(""), func(key []byte, value bool) bool {
		inOrder = append(inOrder, string(key))
		return false
	})

	if len(inOrder) != len(s) {
		t.Fatalf("bad length: %v %v %v", len(inOrder), len(s), r.Len())
	}

	if !slices.Equal(inOrder, s) {
		t.Fatalf("bad order: %v %v", inOrder, s)
	}

	r1 := NewRadixMap[byte, bool]()
	_, ok := r1.Delete([]byte(""))
	if ok {
		t.Fatalf("bad")
	}
	_, ok = r1.ReplaceOrInsert([]byte(""), true)
	if ok {
		t.Fatalf("bad")
	}
	val, ok := r1.Get([]byte(""))
	if !ok || val != true {
		t.Fatalf("bad: %v", val)
	}
	val, ok = r1.Delete([]byte(""))
	if !ok || val != true {
		t.Fatalf("bad: %v", val)
	}

}

func TestEmptyKeyExtend(t *testing.T) {
	r := NewRadixMap[byte, bool]()

	s := []string{"", "A", "AB"}

	for _, ss := range s {
		r.ReplaceOrInsertExtend([]byte(ss), true, nil)
	}

	var inOrder []string

	r.AscendGreaterOrEqual([]byte(""), func(key []byte, value bool) bool {
		inOrder = append(inOrder, string(key))
		return false
	})

	if len(inOrder) != len(s) {
		t.Fatalf("bad length: %v %v %v", len(inOrder), len(s), r.Len())
	}

	if !slices.Equal(inOrder, s) {
		t.Fatalf("bad order: %v %v", inOrder, s)
	}

	r1 := NewRadixMap[byte, bool]()
	_, ok := r1.Delete([]byte(""))
	if ok {
		t.Fatalf("bad")
	}
	_, ok = r1.ReplaceOrInsertExtend([]byte(""), true, nil)
	if ok {
		t.Fatalf("bad")
	}
	val, ok := r1.Get([]byte(""))
	if !ok || val != true {
		t.Fatalf("bad: %v", val)
	}
	val, ok = r1.Delete([]byte(""))
	if !ok || val != true {
		t.Fatalf("bad: %v", val)
	}

}

func TestDelete(t *testing.T) {

	r := NewRadixMap[byte, bool]()

	s := []string{"", "A", "AB"}

	//s := []string{"A", "AB"}

	for _, ss := range s {
		r.ReplaceOrInsert([]byte(ss), true)
	}

	for _, ss := range s {
		_, ok := r.Delete([]byte(ss))

		if !ok {
			t.Fatalf("bad %q", ss)
		}
	}
}

func TestDeletePrefix(t *testing.T) {
	type exp struct {
		inp        []string
		prefix     string
		out        []string
		numDeleted int
	}

	cases := []exp{
		{[]string{"", "A", "AB", "ABC", "R", "S"}, "A", []string{"", "R", "S"}, 3},
		{[]string{"", "A", "AB", "ABC", "R", "S"}, "ABC", []string{"", "A", "AB", "R", "S"}, 1},
		{[]string{"", "A", "AB", "ABC", "R", "S"}, "", []string{}, 6},
		{[]string{"", "A", "AB", "ABC", "R", "S"}, "S", []string{"", "A", "AB", "ABC", "R"}, 1},
		{[]string{"", "A", "AB", "ABC", "R", "S"}, "SS", []string{"", "A", "AB", "ABC", "R", "S"}, 0},
	}

	for _, test := range cases {
		r := NewRadixMap[byte, bool]()
		for _, ss := range test.inp {
			r.ReplaceOrInsert([]byte(ss), true)
		}

		deleted := r.DeletePrefix([]byte(test.prefix))

		if deleted != nil && deleted.len != test.numDeleted {
			t.Fatalf("Bad delete, expected %v to be deleted but got %v", test.numDeleted, deleted)
		}

		out := []string{}

		r.Ascend(func(s []byte, v bool) bool {
			out = append(out, string(s))
			return false
		})

		if !reflect.DeepEqual(out, test.out) {
			t.Fatalf("mis-match: %v %v", out, test.out)
		}
	}
}

func TestLongestPrefix(t *testing.T) {
	r := NewRadixMap[byte, any]()

	keys := []string{
		"",
		"foo",
		"foobar",
		"foobarbaz",
		"foobarbazzip",
		"foozip",
	}
	for _, k := range keys {
		r.ReplaceOrInsert([]byte(k), nil)
	}
	if r.Len() != len(keys) {
		t.Fatalf("bad len: %v %v", r.Len(), len(keys))
	}

	type exp struct {
		inp string
		out string
	}
	cases := []exp{
		{"a", ""},
		{"abc", ""},
		{"fo", ""},
		{"foo", "foo"},
		{"foob", "foo"},
		{"foobar", "foobar"},
		{"foobarba", "foobar"},
		{"foobarbaz", "foobarbaz"},
		{"foobarbazzi", "foobarbaz"},
		{"foobarbazzip", "foobarbazzip"},
		{"foozi", "foo"},
		{"foozip", "foozip"},
		{"foozipzap", "foozip"},
	}
	for _, test := range cases {

		m, _, _ := r.LongestPrefix([]byte(test.inp))

		//if  len(m) != 0 && !ok {
		//	t.Fatalf("no match: %v", test)
		//}
		if string(m) != test.out {
			t.Fatalf("mis-match: %v %v", m, test)
		}
	}
}

func TestWalkPrefix(t *testing.T) {
	r := NewRadixMap[byte, any]()

	keys := []string{
		"foobar",
		"foo/bar/baz",
		"foo/baz/bar",
		"foo/zip/zap",
		"zipzap",
	}
	for _, k := range keys {
		r.ReplaceOrInsert([]byte(k), nil)
	}
	if r.Len() != len(keys) {
		t.Fatalf("bad len: %v %v", r.Len(), len(keys))
	}

	type exp struct {
		inp string
		out []string
	}
	cases := []exp{
		{
			"f",
			[]string{"foobar", "foo/bar/baz", "foo/baz/bar", "foo/zip/zap"},
		},
		{
			"foo",
			[]string{"foobar", "foo/bar/baz", "foo/baz/bar", "foo/zip/zap"},
		},
		{
			"foob",
			[]string{"foobar"},
		},
		{
			"foo/",
			[]string{"foo/bar/baz", "foo/baz/bar", "foo/zip/zap"},
		},
		{
			"foo/b",
			[]string{"foo/bar/baz", "foo/baz/bar"},
		},
		{
			"foo/ba",
			[]string{"foo/bar/baz", "foo/baz/bar"},
		},
		{
			"foo/bar",
			[]string{"foo/bar/baz"},
		},
		{
			"foo/bar/baz",
			[]string{"foo/bar/baz"},
		},
		{
			"foo/bar/bazoo",
			[]string{},
		},
		{
			"z",
			[]string{"zipzap"},
		},
	}

	for _, test := range cases {
		var out []string

		r.AscendPrefix([]byte(test.inp), func(s []byte, v any) bool {
			out = append(out, string(s))
			return false
		})
		sort.Strings(out)
		sort.Strings(test.out)
		if !slices.Equal(out, test.out) {
			t.Fatalf("mis-match: %v %v", out, test.out)
		}
	}
}

func ToMap(r *RadixMap[byte, any]) map[string]any {
	m := make(map[string]any)
	r.Ascend(func(s []byte, v any) bool {
		m[string(s)] = v
		return false
	})
	return m
}

func countBlock(b *block[byte, any]) int {
	count := 1
	for _, c := range b.next {
		if c != nil {
			count += countBlock(c)
		}
	}
	return count
}

func TestWalkMarkDelete(t *testing.T) {
	// Don't modify this test case. This is a test for these specifics keys only. Add keys or modify keys may cause the test to fail.
	r := NewRadixMap[byte, any]()
	r.ReplaceOrInsert([]byte("init0/0"), nil)
	r.ReplaceOrInsert([]byte("init0/1"), nil)
	r.ReplaceOrInsert([]byte("init0/2"), nil)
	r.ReplaceOrInsert([]byte("init0/3"), nil)
	r.ReplaceOrInsert([]byte("init1/0"), nil)
	r.ReplaceOrInsert([]byte("init1/1"), nil)
	r.ReplaceOrInsert([]byte("init1/2"), nil)
	r.ReplaceOrInsert([]byte("init1/3"), nil)
	r.ReplaceOrInsert([]byte("init2"), nil)

	deleteFn := func(s []byte, v any) bool {
		r.MarkDelete(s)
		return false
	}
	r.AscendPrefix([]byte("init1"), deleteFn)

	for _, s := range []string{"init0/0", "init0/1", "init0/2", "init0/3", "init2"} {
		if _, ok := r.Get([]byte(s)); !ok {
			t.Fatalf("expecting to still find %q", s)
		}
	}

	if n := r.Len(); n != 5 {
		t.Fatalf("expected to find exactly 5 nodes, instead found %d: %v", n, ToMap(r))
	}

	if countBlock(r.root) != 13 {
		t.Fatalf("expected to find exactly 13 blocks, instead found %d", countBlock(r.root))
	}

	r.Shrink()

	r.Ascend(func(key []byte, v any) bool {
		return false
	})

	if countBlock(r.root) != 8 {
		t.Fatalf("expected to find exactly 8 blocks, instead found %d", countBlock(r.root))
	}

	r.Ascend(func(key []byte, v any) bool {
		return false
	})

	r.Ascend(deleteFn)

	if n := r.Len(); n != 0 {
		t.Fatalf("expected to find exactly 0 nodes, instead found %d: %v", n, ToMap(r))
	}

	if countBlock(r.root) != 8 {
		t.Fatalf("expected to find exactly 8 blocks, instead found %d", countBlock(r.root))
	}
	r.Shrink()
	if countBlock(r.root) != 1 {
		t.Fatalf("expected to find exactly 1 block, instead found %d", countBlock(r.root))
	}

}

func panics(f func()) (b bool) {
	defer func() {
		if x := recover(); x != nil {
			b = true
		}
	}()
	f()
	return false
}

func TestBug_1(t *testing.T) {
	// test for a bug that a panic is triggered when the key is end with fully matched label in a non-leaf node.
	if panics(func() {
		r := NewRadixMap[byte, int]()
		r.ReplaceOrInsert([]byte("init0/0"), 0)
		r.ReplaceOrInsert([]byte("init0/1"), 0)
		r.ReplaceOrInsert([]byte("init0/2"), 0)
		r.ReplaceOrInsert([]byte("init0/3"), 0)
		r.ReplaceOrInsert([]byte("init0/"), 0)

		r.ReplaceOrInsertExtend([]byte("init1/0"), 0, func(v int) int { return v + 1 })
		r.ReplaceOrInsertExtend([]byte("init1/1"), 0, func(v int) int { return v + 1 })
		r.ReplaceOrInsertExtend([]byte("init1/2"), 0, func(v int) int { return v + 1 })
		r.ReplaceOrInsertExtend([]byte("init1/3"), 0, func(v int) int { return v + 1 })
		r.ReplaceOrInsertExtend([]byte("init0/"), 0, func(v int) int { return v + 1 })
	}) {
		t.Errorf("Should not panic")
	}

}
