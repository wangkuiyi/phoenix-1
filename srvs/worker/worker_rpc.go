package main

import (
	"bufio"
	"encoding/gob"
	"math/rand"

	"github.com/huichen/sego"
	"github.com/wangkuiyi/fs"
	"github.com/wangkuiyi/phoenix/algo"
	"github.com/wangkuiyi/phoenix/srvs"
)

// Worker is the RPC service implementation.
type Worker struct {
	sgmt  *sego.Segmenter
	vocab *algo.Vocab
	vshdr *algo.VSharder
	rng   *rand.Rand
	cfg   *srvs.Config
}

func (w *Worker) Initialize(shard_filename string, _ *int) error {
	in, e := fs.Open(shard_filename)
	if e != nil {
		return e
	}
	defer in.Close()

	out, e := fs.Create(w.cfg.IterDir(0))
	if e != nil {
		return e
	}
	defer out.Close()

	scanner := bufio.NewScanner(in)
	encoder := gob.NewEncoder(out)
	for scanner.Scan() {
		d := algo.NewDocument(scanner.Text(), w.sgmt, w.vocab, w.vshdr, w.rng, w.cfg.Topics)
		if e := encoder.Encode(d); e != nil {
			return e
		}
	}
	return scanner.Err()
}
