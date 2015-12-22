package main

import (
	"bufio"
	"encoding/gob"
	"fmt"

	"github.com/wangkuiyi/fs"
	"github.com/wangkuiyi/phoenix/algo"
	"github.com/wangkuiyi/phoenix/srvs"
	"github.com/wangkuiyi/sego"
)

// Worker is the RPC service implementation.
type WorkerRPC struct {
	addr  string // Worker address. Also worker ID.
	sgmt  *sego.Segmenter
	vocab *algo.Vocab
	vshdr *algo.VSharder
	cfg   *srvs.Config
}

func (w *WorkerRPC) Initialize(shard_filename string, _ *int) error {
	if w.sgmt == nil {
		w.sgmt = new(sego.Segmenter)
		w.sgmt.LoadDictionary(w.cfg.Segmenter)
	}

	rng := srvs.NewRand(shard_filename)

	in, e := fs.Open(shard_filename)
	if e != nil {
		return fmt.Errorf("%v.Initialize(%v): %v", w.addr, shard_filename, e)
	}
	defer in.Close()

	out, e := fs.Create(w.cfg.IterDir(0))
	if e != nil {
		return fmt.Errorf("%v.Initialize(%v): %v", w.addr, shard_filename, e)
	}
	defer out.Close()

	scanner := bufio.NewScanner(in)
	encoder := gob.NewEncoder(out)
	for scanner.Scan() {
		d := algo.NewDocument(scanner.Text(), w.sgmt, w.vocab, w.vshdr, rng, w.cfg.Topics)
		if e := encoder.Encode(d); e != nil {
			return fmt.Errorf("%v.Initialize(%v): %v", w.addr, shard_filename, e)
		}
	}
	return fmt.Errorf("%v.Initialize(%v): %v", w.addr, shard_filename, scanner.Err())
}
