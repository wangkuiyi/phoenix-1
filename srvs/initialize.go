package srvs

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"log"
	"path"

	"github.com/wangkuiyi/fs"
	"github.com/wangkuiyi/parallel"
	"github.com/wangkuiyi/phoenix/algo"
)

type InitializeArg struct {
	Shard       string
	Aggregators []*RPC
}

func (m *Master) Initialize() {
	log.Println("Initialization ...")
	defer log.Println("Initialization done")

	if fis, e := fs.ReadDir(m.cfg.CorpusDir); e != nil {
		log.Panicf("Failed listing corpus shards in %s: %v", m.cfg.CorpusDir, e)
	} else {
		// TODO(y): create a temp dir and rename after finishing.
		if e := fs.Mkdir(m.cfg.IterPath(0)); e != nil {
			log.Panicf("Failed creating output directory %s: %v", m.cfg.IterPath(0), e)
		}

		ch := make(chan string)
		go parallel.For(0, len(m.workers), 1, func(i int) {
			for fn := range ch {
				var dumb int
				if e := m.workers[i].Call("Worker.Initialize", &InitializeArg{fn, m.aggregators}, &dumb); e != nil {
					log.Panicf("Worker %s failed initializing shard %s: %v", m.workers[i].Addr, fn, e)
				}
			}
		})
		for _, fi := range fis {
			if !fi.IsDir() {
				ch <- fi.Name()
			}
		}
		close(ch)
	}
}

func (w *Worker) Initialize(arg *InitializeArg, _ *int) error {
	log.Printf("Worker(%s).Initialize(%s) ...", w.addr, arg.Shard)

	if e := GuaranteeSegmenter(&w.sgmt, w.cfg.Segmenter); e != nil {
		return fmt.Errorf("Worker %s cannot load segmenter dictionary %s: %v", w.addr, w.cfg.Segmenter, e)
	}
	if e := GuaranteeVocabSharder(&w.vocab, &w.vshdr, w.cfg.Vocab, w.cfg.VShards); e != nil {
		return fmt.Errorf("Aggregator %v failed to build vocab and vsharder from %v: %v", w.addr, w.cfg.Vocab, e)
	}
	model := algo.NewModel(true /*dense*/, w.vocab, nil /*global model*/, 0, w.cfg.Topics)

	rng := NewRand(arg.Shard) //TODO(y): May need more randomness here.

	in, e := fs.Open(path.Join(w.cfg.CorpusDir, arg.Shard))
	if e != nil {
		return fmt.Errorf("Worker(%v).Initialize(%v): %v", w.addr, arg.Shard, e)
	}
	defer in.Close()

	out, e := fs.Create(path.Join(w.cfg.IterPath(0), arg.Shard))
	if e != nil {
		return fmt.Errorf("Worker(%v).Initialize(%v): %v", w.addr, arg.Shard, e)
	}
	defer out.Close()

	scanner := bufio.NewScanner(in)
	encoder := gob.NewEncoder(out)
	for scanner.Scan() {
		d := algo.NewDocument(scanner.Text(), w.sgmt, w.vocab, w.vshdr, rng, w.cfg.Topics, model)
		if e := encoder.Encode(d); e != nil {
			return fmt.Errorf("%v.Initialize(%v): %v", w.addr, arg.Shard, e)
		}
	}
	if scanner.Err() != nil {
		return fmt.Errorf("%v.Initialize(%v): %v", w.addr, arg.Shard, scanner.Err())
	}

	log.Printf("Worker(%s).Initialize(%s) done", w.addr, arg.Shard)
	return nil
}
