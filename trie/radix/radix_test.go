package radix

import (
	crand "crypto/rand"
	"fmt"
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

func inOrderByteTraversal(n *node[byte, int]) []string {
	if n == nil {
		return nil
	}

	var ret []string

	var inOrderRecursive func(n *node[byte, int], prefix []byte)

	inOrderRecursive = func(n *node[byte, int], prefix []byte) {

		if n.lastElem {
			ret = append(ret, string(prefix))
		}

		for _, e := range n.edges {
			inOrderRecursive(e.node, append(prefix, e.label...))
		}
	}

	inOrderRecursive(n, []byte{})

	return ret
}

func TestNewRadixTrieMap(t *testing.T) {
	var min, max string
	inp := make(map[string]int)
	for i := 0; i < 1000; i++ {
		gen := generateUUID()
		inp[gen] = i
		if gen < min || i == 0 {
			min = gen
		}
		if gen > max || i == 0 {
			max = gen
		}
	}

	r := NewRadixTrieMap[byte, int]()

	for k, v := range inp {
		r.ReplaceOrInsert([]byte(k), v)
	}

	if r.Len() != len(inp) {
		t.Errorf("Got %v expected %v", r.Len(), len(inp))
	}

	//r.Insert([]byte("romane"), 2)
	//fmt.Println("in order traverse: ", inOrderByteTraversal(r.root))
	//
	//r.Insert([]byte("romanus"), 3)
	//
	//fmt.Println("in order traverse: ", inOrderByteTraversal(r.root))
	//
	//if r.Len() != len(inp)+2 {
	//	t.Errorf("Got %v expected %v", r.Len(), len(inp)+2)
	//}

}
