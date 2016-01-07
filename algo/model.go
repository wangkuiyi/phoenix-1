package algo

import (
	"log"

	"github.com/wangkuiyi/phoenix/algo/hist"
)

type Model struct {
	Vocab       *Vocab
	VShdr       *VSharder // If VShdr is nil, this is a global model.
	VShard      int
	DenseHists  []hist.Dense // either use Dense or Sparse
	SparseHists []hist.Sparse
	Global      hist.Dense64
}

// Create either a dense model (for fast accessing) or a sparse model (for saving memory).
func NewModel(dense bool, vocab *Vocab, vshdr *VSharder, vshard int, topics int) *Model {
	if vshdr != nil && (vshard < 0 || vshard >= vshdr.Num()) {
		log.Panicf("vshard (%d) out of range [0, %d)", vshard, vshdr.Num())
	}
	m := &Model{
		Vocab:  vocab,
		VShdr:  vshdr,
		VShard: vshard,
		Global: make(hist.Dense64, topics)}

	tokens := len(vocab.Ids)
	if vshdr != nil {
		tokens = vshdr.Size(vshard)
	}

	if dense {
		m.DenseHists = make([]hist.Dense, tokens)
	} else {
		m.SparseHists = make([]hist.Sparse, tokens)
	}
	return m
}

func (m *Model) Topics() int {
	return len(m.Global)
}

// If a token is in this model, return true and VShard-local token Id.
func (m *Model) In(token int) (bool, int) {
	if m.VShdr == nil {
		return true, token
	}
	offset := token - m.VShdr.Begin(m.VShard)
	return 0 <= offset && offset < m.VShdr.Size(m.VShard), offset
}

func (m *Model) Inc(token, topic int, count int32) {
	// NOTE: No panic for tokens out of VShard.
	if in, off := m.In(token); in {
		if m.DenseHists != nil { // This is a dense model.
			if m.DenseHists[off] == nil {
				m.DenseHists[off] = make(hist.Dense, m.Topics())
			}
			m.DenseHists[off][topic] += count
			m.Global[topic] += int64(count)
		} else if m.SparseHists != nil {
			if m.SparseHists[off] == nil {
				m.SparseHists[off] = make(hist.Sparse)
			}
			m.SparseHists[off][topic] += count
			m.Global[topic] += int64(count)
		} else {
			log.Panicf("Both Model.DenseHists and Model.SparseHists are nil.")
		}
	}
}

func (m *Model) Dense(token int) hist.Dense {
	if in, off := m.In(token); in {
		return m.DenseHists[off]
	}
	return nil
}

func (m *Model) Sparse(token int) hist.Sparse {
	if in, off := m.In(token); in {
		return m.SparseHists[off]
	}
	return nil
}
