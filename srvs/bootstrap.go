package srvs

import (
	"log"

	"github.com/wangkuiyi/fs"
	"github.com/wangkuiyi/parallel"
)

type BootstrapArg struct {
	Iter   int
	VShard int
}

// Bootstrap notifies the first Config.VShards of all registered
// aggregators and all workers to bootstrap.
func (m *Master) bootstrap(iter int) {
	log.Println("Bootstrapping...")

	if iter < 0 {
		fs.Mkdir(m.cfg.BaseDir)
	}

	if e := GuaranteeConfig(m.cfg); e != nil {
		log.Panic(e)
	}

	if iter >= 0 {
		parallel.For(0, len(m.workers), 1, func(i int) {
			arg := BootstrapArg{iter, i % m.cfg.VShards}
			if e := m.workers[i].Call("Worker.Bootstrap", &arg, nil); e != nil {
				log.Panicf("Worker(%s).Bootstrap(%v) failed: %v", m.workers[i].Addr, arg, e)
			}
		})
	}

	parallel.For(0, m.cfg.VShards, 1, func(i int) {
		arg := BootstrapArg{iter, i}
		if e := m.aggregators[i].Call("Aggregator.Bootstrap", &arg, nil); e != nil {
			log.Panicf("Aggregator(%s).Bootstrap(%v) failed: %v", m.aggregators[i].Addr, iter, e)
		}
	})

	log.Println("Bootstrap done")
}

// Bootstrap loads existing model if arg.Iter >= 0.
func (w *Worker) Bootstrap(arg *BootstrapArg, _ *int) error {
	log.Printf("Worker(%s).Bootstrap(%+v) ...", w.addr, arg)
	e := GuaranteeModel(&w.model, &w.vocab, &w.vshdr, &w.cfg, arg)
	log.Printf("Worker(%s).Boostrap(%+v) done", w.addr, arg)
	return e
}

// Bootstrap loads existing model if arg.Iter >= 0, or create empty models
func (a *Aggregator) Bootstrap(arg *BootstrapArg, _ *int) error {
	log.Printf("Aggregator(%s).Bootstrap(%+v) ...", a.addr, arg)
	e := GuaranteeModel(&a.model, &a.vocab, &a.vshdr, &a.cfg, arg)
	log.Printf("Aggregator(%s).Boostrap(%+v) done", a.addr, arg)
	return e
}
