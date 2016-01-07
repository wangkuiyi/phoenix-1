package srvs

import (
	"encoding/gob"
	"fmt"
	"log"

	"github.com/wangkuiyi/fs"
	"github.com/wangkuiyi/parallel"
	"github.com/wangkuiyi/phoenix/algo"
)

type BootstrapArg struct {
	Iter   int
	VShard int
}

// Bootstrap notifies the first Config.VShards of all registered
// aggregators to load an existing model shard or initialize a new one.
func (m *Master) Bootstrap(iter int) {
	log.Println("Bootstrapping...")

	if iter < 0 {
		fs.Mkdir(m.cfg.BaseDir)
	}

	if e := GuaranteeConfig(m.cfg); e != nil {
		log.Panic(e)
	}

	parallel.For(0, m.cfg.VShards, 1, func(i int) {
		var dumb int
		if e := m.aggregators[i].Call("Aggregator.Bootstrap", &BootstrapArg{iter, i}, &dumb); e != nil {
			log.Panicf("Aggregator(%s).Bootstrap(%v) failed: %v", m.aggregators[i].Addr, iter, e)
		}
	})

	log.Println("Bootstrap done")
}

func (a *Aggregator) Bootstrap(arg *BootstrapArg, _ *int) error {
	log.Printf("Aggregator(%s).Bootstrap(%+v) ...", a.addr, arg)

	if arg.Iter < 0 { // Initialziation is not yet completed.
		if e := GuaranteeVocabSharder(&a.vocab, &a.vshdr, a.cfg.Vocab, a.cfg.VShards); e != nil {
			return fmt.Errorf("Aggregator %v failed to build vocab and vsharder from %v: %v", a.addr, a.cfg.Vocab, e)
		}
		a.model = algo.NewModel(true, a.vocab, a.vshdr, arg.VShard, a.cfg.Topics)
	} else {
		if model, e := fs.Open(a.cfg.VShardPath(arg.Iter, arg.VShard)); e != nil {
			return fmt.Errorf("Aggregator %v failed to open vshard %+v: %v", a.addr, arg, e)
		} else {
			defer model.Close()
			a.model = &algo.Model{}
			if e := gob.NewDecoder(model).Decode(a.model); e != nil {
				return fmt.Errorf("Aggregator %v failed to decode model vshard %+v: %v", a.addr, arg, e)
			}
			a.vocab = a.model.Vocab
			a.vshdr = a.model.VShdr
		}
	}

	log.Printf("Aggregator(%s).Boostrap(%+v) done", a.addr, arg)
	return nil
}
