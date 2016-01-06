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
	"github.com/wangkuiyi/sego"
)

func (m *Master) Initialize() {
	if fis, e := fs.ReadDir(m.CorpusDir); e != nil {
		log.Panicf("Failed listing corpus shards in %s: %v", m.CorpusDir, e)
	} else {
		ch := make(chan string)
		parallel.For(0, len(m.workers), 1, func(i int) {
			for fn := range ch {
				var dumb int
				if e := m.workers[i].Call("Worker.Initialize", fn, &dumb); e != nil {
					log.Panicf("Worker %s failed initializing shard %s: %v", m.workers[i].Addr, fn, e)
				}
			}
		})
		for _, fi := range fis {
			if !fi.IsDir() {
				ch <- path.Join(m.CorpusDir, fi.Name())
			}
		}
		close(ch)
	}
}

func (w *Worker) Initialize(shard_filename string, _ *int) error {
	if w.sgmt == nil {
		w.sgmt = new(sego.Segmenter)
		w.sgmt.LoadDictionary(w.Segmenter)
	}

	rng := NewRand(shard_filename)

	in, e := fs.Open(shard_filename)
	if e != nil {
		return fmt.Errorf("%v.Initialize(%v): %v", w.addr, shard_filename, e)
	}
	defer in.Close()

	out, e := fs.Create(w.IterPath(0))
	if e != nil {
		return fmt.Errorf("%v.Initialize(%v): %v", w.addr, shard_filename, e)
	}
	defer out.Close()

	scanner := bufio.NewScanner(in)
	encoder := gob.NewEncoder(out)
	for scanner.Scan() {
		d := algo.NewDocument(scanner.Text(), w.sgmt, w.vocab, w.vshdr, rng, w.Topics)
		if e := encoder.Encode(d); e != nil {
			return fmt.Errorf("%v.Initialize(%v): %v", w.addr, shard_filename, e)
		}
	}
	return fmt.Errorf("%v.Initialize(%v): %v", w.addr, shard_filename, scanner.Err())
}
