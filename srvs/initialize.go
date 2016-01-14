package srvs

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"log"
	"path"
	"time"

	"github.com/wangkuiyi/fs"
	"github.com/wangkuiyi/parallel"
	"github.com/wangkuiyi/phoenix/algo"
	"github.com/wangkuiyi/phoenix/algo/hist"
)

type InitializeArg struct {
	Shard       string
	Aggregators []*RPC
}

func (m *Master) initialize() {
	log.Println("Initialization ...")
	start := time.Now()

	tmpDir := path.Join(m.cfg.BaseDir, "current")
	if e := fs.Mkdir(tmpDir); e != nil {
		log.Panicf("Failed to create initialization output directory %s: %v", tmpDir, e)
	}

	modelDir := path.Join(tmpDir, "model")
	if e := fs.Mkdir(modelDir); e != nil {
		log.Panicf("Failed to create initialziation model directory %s : %v", modelDir, e)
	}

	ch := make(chan string)
	go func() {
		for _, fn := range m.corpusShards {
			ch <- fn
		}
		close(ch)
	}()
	parallel.For(0, len(m.workers), 1, func(i int) {
		for fn := range ch {
			var dumb int
			if e := m.workers[i].Call("Worker.Initialize", &InitializeArg{fn, m.aggregators}, &dumb); e != nil {
				log.Panicf("Worker %s failed initializing shard %s: %v", m.workers[i].Addr, fn, e)
			}
		}
	})

	parallel.For(0, m.cfg.VShards, 1, func(v int) {
		var dumb int
		if e := m.aggregators[v].Call("Aggregator.Save", tmpDir, &dumb); e != nil {
			log.Panicf("Aggregator %s failed to save model: %v", m.aggregators[v].Addr, e)
		}
	})

	fs.Rename(tmpDir, m.cfg.IterPath(0))
	log.Printf("Initialization done in %s", time.Since(start))
}

func (w *Worker) Initialize(arg *InitializeArg, _ *int) error {
	log.Printf("Worker(%s).Initialize(%s) ...", w.addr, arg.Shard)

	if e := GuaranteeSegmenter(&w.sgmt, w.cfg.Segmenter); e != nil {
		return fmt.Errorf("Worker %s cannot load segmenter dictionary %s: %v", w.addr, w.cfg.Segmenter, e)
	}
	if e := GuaranteeVocabSharder(&w.vocab, &w.vshdr, w.cfg.Vocab, w.cfg.VShards); e != nil {
		return fmt.Errorf("Aggregator %v failed to build vocab and vsharder from %v: %v", w.addr, w.cfg.Vocab, e)
	}

	vshards := make([]*algo.Model, w.cfg.VShards)
	for v := range vshards {
		vshards[v] = algo.NewModel(true /*dense*/, w.vocab, w.vshdr, v, w.cfg.Topics) // TODO(y): Check real sparsity and switch between dense/sparse.
	}

	rng := NewRand(arg.Shard) //TODO(y): May need more randomness here.

	in, e := fs.Open(path.Join(w.cfg.CorpusDir, arg.Shard))
	if e != nil {
		return fmt.Errorf("Worker(%v).Initialize(%v): %v", w.addr, arg.Shard, e)
	}
	defer in.Close()

	out, e := fs.Create(path.Join(w.cfg.BaseDir, "current", arg.Shard))
	if e != nil {
		return fmt.Errorf("Worker(%v).Initialize(%v): %v", w.addr, arg.Shard, e)
	}
	defer out.Close()

	scanner := bufio.NewScanner(in)
	encoder := gob.NewEncoder(out)
	for scanner.Scan() {
		d := algo.NewDocument(scanner.Text(), w.sgmt, w.vocab, w.vshdr, rng, w.cfg.Topics, vshards)
		if e := encoder.Encode(d); e != nil {
			return fmt.Errorf("%v.Initialize(%v): %v", w.addr, arg.Shard, e)
		}
	}
	if scanner.Err() != nil {
		return fmt.Errorf("%v.Initialize(%v) scanner error: %v", w.addr, arg.Shard, scanner.Err())
	}

	parallel.For(0, w.cfg.VShards, 1, func(v int) error {
		if e := arg.Aggregators[v].Dial(); e != nil {
			return e
		}
		if e := arg.Aggregators[v].Call("Aggregator.Aggregate", vshards[v].DenseHists, nil); e != nil {
			return e
		}
		return nil
	})

	log.Printf("Worker(%s).Initialize(%s) done", w.addr, arg.Shard)
	return nil
}

func (a *Aggregator) Aggregate(diff []hist.Dense, _ *int) error {
	log.Println("Aggregator.Aggregate ...")

	if a.model.DenseHists == nil || a.model.SparseHists != nil {
		return fmt.Errorf("Aggregator.Aggregate(diff): model must be dense")
	}
	if len(a.model.DenseHists) != len(diff) {
		return fmt.Errorf("Aggregator.Aggregate(diff): model vshard size (%d) != diff size(%d)", len(a.model.DenseHists), len(diff))
	}
	for i := range diff {
		for t := 0; t < len(diff[i]); t++ {
			if a.model.DenseHists[i] == nil {
				a.model.DenseHists[i] = make(hist.Dense, a.model.Topics())
			}
			a.model.DenseHists[i][t] += diff[i][t]
			a.model.Global[t] += int64(diff[i][t])
		}
	}

	log.Println("Aggregator.Aggregate done")
	return nil
}

func (a *Aggregator) Save(iterDir string, _ *int) error {
	log.Printf("Aggregator.Save(%s) ...", iterDir)

	vshard := path.Join(iterDir, VShardName(a.model.VShard, a.cfg.VShards))
	out, e := fs.Create(vshard)
	if e != nil {
		return fmt.Errorf("Aggregator.Save() cannot create vshard file %v: %v", vshard, e)
	}
	defer out.Close()
	e = gob.NewEncoder(out).Encode(a.model)

	log.Printf("Aggregator.Save(%s) done", iterDir)
	return e
}
