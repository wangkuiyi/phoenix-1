// hist defines histogram data structures used in package algo.  It
// assumes that (1) topic-word co-occurrences are int32, and (2)
// global topic counts are int64.
package hist

// Sparse is for buffering hist diffs.
type Sparse map[int]int32

// Dense is for topic-histogram of words in models.
type Dense []int32

// Dense64 is for the global topic-histogram.
type Dense64 []int64

// Ordered is used in Document to support fast Gibbs sampling.
type Ordered struct {
	Topics []int32
	Tokens []int32
}
