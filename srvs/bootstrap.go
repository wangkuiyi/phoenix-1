package srvs

import (
	"encoding/gob"
	"fmt"
	"log"
	"path"

	"github.com/wangkuiyi/fs"
	"github.com/wangkuiyi/parallel"
	"github.com/wangkuiyi/phoenix/algo"
)

// Bootstrap notifies the first Config.VShards aggregators to load
// an existing model shard, or to initialize a new model shard.
func (m *Master) Bootstrap(iter int) {
	if iter < 0 {
		fs.Mkdir(m.BaseDir)
	}

	reordered := make([]*RPC, len(m.aggregators))
	parallel.For(0, m.VShards, 1, func(i int) {
		var realVShard int
		if e := m.aggregators[i].Call("Aggregator.Bootstrap", m.VShardPath(iter, i), &realVShard); e != nil {
			log.Panicf("Aggregator(%s).Bootstrap(%v) failed: %v", m.aggregators[i].Addr, iter, e)
		} else {
			reordered[realVShard] = m.aggregators[i]
		}
	})
	m.aggregators = reordered
}

// Bootstrap creates a new model shard if vshardPath is invalid, or
// load an existing model shard.  Anyway, it returns the real vshard
// sequence number, so that master can order RPC clients of all
// aggregators.  This re-order will be sent to workers during
// initialization and Gibbs sampling.
func (a *Aggregator) Bootstrap(vshardPath string, realVShard *int) error {
	tokens, e := fs.Open(a.Vocab)
	if e != nil {
		return fmt.Errorf("Aggregator %v failed to open vocab file: %v", a.addr, e)
	}
	defer tokens.Close()

	a.vocab, a.vshdr, e = algo.BuildVocabAndVSharder(tokens, a.VShards, true)
	if e != nil {
		return fmt.Errorf("Aggregator %v failed to build vocab and vsharder: %v", a.addr, e)
	}

	vshard, _, e := VShardFromName(path.Base(vshardPath))
	if e != nil {
		return fmt.Errorf("Aggregator %v failed to parse vshard from path %v: %v", a.addr, vshardPath, e)
	}

	iter, e := IterFromDir(path.Base(path.Dir(vshardPath)))
	if e != nil || iter < 0 { // e != nil when the most recent iteration is negative
		a.model = algo.NewModel(true, a.vocab, a.vshdr, vshard, a.Topics)
		*realVShard = vshard
	} else {
		if model, e := fs.Open(vshardPath); e != nil {
			return fmt.Errorf("Aggregator %v failed to open vshard file %v: %v", a.addr, vshardPath, e)
		} else {
			defer model.Close()
			a.model = &algo.Model{}
			if e := gob.NewDecoder(model).Decode(a.model); e != nil {
				return fmt.Errorf("Aggregator %v failed to decode model from %v: %v", a.addr, vshardPath, e)
			}
			*realVShard = a.model.VShard
		}
	}

	return nil
}
