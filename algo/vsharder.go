package algo

import (
	"bufio"
	"container/heap"
	"fmt"
	"io"
	"log"
	"sort"
	"strings"
)

// Vocab is a mapping between token and token ID.
type Vocab struct {
	Ids    map[string]int
	Tokens []string
	// TODO(y): Add Freq[]float64 so we can print shard freqs.
}

// VSharder is a mapping between token ID and multi-level of V-shards.
type VSharder struct {
	Breakpoints []int // Accumulative bucket sizes.
}

// If Vocab contains token, returns its Id; otherwise returns -1.
func (v *Vocab) Id(token string) int {
	if id, ok := v.Ids[token]; ok {
		return id
	}
	return -1
}

// If Vocab contains token, returns token; otherwise panics.
func (v *Vocab) Token(id int) string {
	if id < 0 || id >= len(v.Tokens) {
		log.Panicf("Vocab.Token: id (%d) is out of range [%d, %d)", id, 0, len(v.Tokens))
	}
	return v.Tokens[id]
}

// Use binary search to find locate shard of a token.
func (v *VSharder) Shard(token int) int {
	if token < 0 {
		return -1
	}
	return sort.Search(len(v.Breakpoints), func(i int) bool { return v.Breakpoints[i] > token })
}

// Returns the real number of shards, after filtering out singular shards.
func (v *VSharder) Num() int {
	return len(v.Breakpoints)
}

func BuildVocabAndVSharder(tokenFreqList io.Reader, numVShards int, delUnbalanced bool) (*Vocab, *VSharder, error) {
	h, e := buildBalancedVShards(tokenFreqList, numVShards)
	if e != nil {
		return nil, nil, e
	}

	v := &Vocab{
		Ids:    make(map[string]int),
		Tokens: make([]string, 0)}
	s := &VSharder{
		Breakpoints: make([]int, 0, len(h))}

	μ, _ := h.variance()
	a := 0
	for _, shard := range h {
		if !delUnbalanced || len(shard.tokens) > 1 || shard.weight <= μ { // NOTE: ignore singular and unbalaned buckets.
			for _, token := range shard.tokens {
				v.Tokens = append(v.Tokens, token)
				v.Ids[token] = len(v.Tokens) - 1
			}
			a += len(shard.tokens)
			s.Breakpoints = append(s.Breakpoints, a)
		}
	}
	return v, s, nil
}

func buildBalancedVShards(r io.Reader, n int) (vshardHeap, error) {
	h := make(vshardHeap, n)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		var freq float64
		var token string
		if _, e := fmt.Fscanf(strings.NewReader(scanner.Text()), "%f %s", &freq, &token); e != nil {
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

func (h vshardHeap) variance() (float64, float64) {
	s := 0.0
	for _, v := range h {
		s += v.weight
	}
	mean := s / float64(len(h))

	s = 0.0
	for _, v := range h {
		s += (v.weight - mean) * (v.weight - mean)
	}
	return mean, s / float64(len(h))
}
