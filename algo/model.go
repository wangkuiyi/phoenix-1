package algo

import (
	"log"

	"github.com/wangkuiyi/phoenix/algo/hist"
)

type Model struct {
	Vocab       *Vocab
	VSharder    *VSharder
	VShard      int
	DenseHists  []hist.Dense // either use Dense or Sparse
	SparseHists []hist.Sparse
	Global      hist.Dense64
}

// Create either a dense model (for fast accessing) or a sparse model (for saving memory).
func NewModel(dense bool, vocab *Vocab, vshdr *VSharder, vshard int, topics int) *Model {
	if vshard < 0 || vshard >= vshdr.Num() {
		log.Panicf("vshard (%d) out of range [0, %d)", vshard, vshdr.Num())
	}
	m := &Model{
		Vocab:    vocab,
		VSharder: vshdr,
		VShard:   vshard,
		Global:   make(hist.Dense64, topics)}
	if dense {
		m.DenseHists = make([]hist.Dense, vshdr.Size(vshard))
	} else {
		m.SparseHists = make([]hist.Sparse, vshdr.Size(vshard))
	}
	return m
}

func (m *Model) Topics() int {
	return len(m.Global)
}

// If a token is in this VShard, return true and VShard-local token Id.
func (m *Model) In(token int) (bool, int) {
	offset := token - m.VSharder.Begin(m.VShard)
	return 0 <= offset && offset < m.VSharder.Size(m.VShard), offset
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
