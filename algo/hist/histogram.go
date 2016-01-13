// hist defines histogram data structures used in package algo.  It
// assumes that (1) topic-word co-occurrences are int32, and (2)
// global topic counts are int64.
package hist

import (
	"bytes"
	"fmt"
	"sort"
)

// Sparse is for buffering hist diffs.
type Sparse map[int]int32

// Dense is for topic-histogram of words in models.
type Dense []int32

// Dense64 is for the global topic-histogram.
type Dense64 []int64

// Ordered is used in Document to support fast Gibbs sampling.  It is
// also used to format model into human readable format.
type Ordered struct {
	Topics []int32
	Counts []int32
}

func OrderedFromDense(d Dense) Ordered {
	o := Ordered{
		Topics: make([]int32, 0, len(d)),
		Counts: make([]int32, 0, len(d))}
	for t, c := range d {
		if c > 0 {
			o.Topics = append(o.Topics, int32(t))
			o.Counts = append(o.Counts, c)
		}
	}
	sort.Sort(&o)
	return o
}

func (o *Ordered) Len() int           { return len(o.Topics) }
func (o *Ordered) Less(i, j int) bool { return o.Counts[i] > o.Counts[j] } // Decensing order of counts.
func (o *Ordered) Swap(i, j int) {
	o.Topics[i], o.Topics[j], o.Counts[i], o.Counts[j] = o.Topics[j], o.Topics[i], o.Counts[j], o.Counts[i]
}

func (o Ordered) String() string {
	var buf bytes.Buffer
	for i := 0; i < 10 && i < len(o.Topics); i++ { // Print at most top 10 topics.
		fmt.Fprintf(&buf, "%d:%d ", o.Topics[i], o.Counts[i])
	}
	return buf.String()
}
