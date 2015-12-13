package algo

import (
	"bufio"
	"container/heap"
	"fmt"
	"io"
	"strings"
)

// Vocab is a mapping between token and token ID.
type Vocab struct {
	Token2Id map[string]int
	Id2Token []string
}

// VSharder is a mapping between token ID and multi-level of V-shards.
type VSharder struct {
	Breakpoints []int // breakpoints of the most detailed level.
	NumShards   []int // number of V-shards at each level.
}

// func BuildVocabAndVSharder(tokenFreqList io.Reader, numVShards ...int) (*Vocab, *VSharder, error) {
// 	tfl, e := loadTokenFreqList(tokenFreqList)
// 	if e != nil {
// 		return nil, nil, e
// 	}

// 	return buildVocab(tfl), buildVShards(tfl, numVShards), e
// }

func buildBalancedVShards(r io.Reader, n int) ([]vshard, error) {
	h := make(vshardHeap, n)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		var freq float64
		var token string
		if _, e := fmt.Fscanf(strings.NewReader(scanner.Text()), "%f %v", &freq, &token); e != nil {
			return nil, e
		}
		h[0].tokens = append(h[0].tokens, token)
		h[0].weight += freq
		heap.Fix(h, 0)
	}
	return h, scanner.Err()
}

// vshardHeap is a min-heap -- the vshard with the least total token frequency is on top.
type vshardHeap []vshard
type vshard struct {
	tokens []string
	weight float64
}

func (h vshardHeap) Len() int {
	return len(h)
}
func (h vshardHeap) Less(i, j int) bool {
	return h[i].weight < h[j].weight
}
func (h vshardHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}
func (h vshardHeap) Push(x interface{}) {
	panic("vshardHeap.Push not implemented")
}
func (h vshardHeap) Pop() interface{} {
	panic("vshardHeap.Pop not implemented")
}

// TODO(y): Add vshardHeap.variance.

func prod(s []int) int {
	r := 1
	for _, v := range s {
		r *= v
	}
	return r
}
