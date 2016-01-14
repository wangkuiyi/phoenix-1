package srvs

import (
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"path"
	"time"

	"github.com/wangkuiyi/fs"
	"github.com/wangkuiyi/parallel"
	"github.com/wangkuiyi/phoenix/algo"
)

func (m *Master) gibbs(iter int) {
	log.Printf("Gibbs %d ...", iter)
	start := time.Now()

	tmpDir := path.Join(m.cfg.BaseDir, "current")
	if e := fs.Mkdir(tmpDir); e != nil {
		log.Panicf("Failed to create initialization output directory %s: %v", tmpDir, e)
	}

	parallel.For(0, len(m.workers), 1, func(w int) {
		// TODO(y): Implement Worker.DownloadModel
		if e := m.workers[w].Call("Worker.DownloadModel", m.aggregators[w%m.cfg.VShards].Addr, nil); e != nil {
			log.Panicf("Worker[%d].DownloadModel(%d) failed: %v", w, w%m.cfg.VShards, e)
		}
	})

	ch := make(chan []string)
	go func() {
		for s := 0; s < len(m.corpusShards); s += m.cfg.VShards {
			// Note: GuaranteeCorpusShardList makes sure that len(m.corpusShards) dividable by VShards.
			ch <- m.corpusShards[s : s+m.cfg.VShards]
		}
		close(ch)
	}()
	parallel.For(0, len(m.workers), m.cfg.VShards, func(w int) {
		for seg := range ch {
			m.sampleSegment(seg, m.workers[w:w+m.cfg.VShards], m.aggregators, iter)
		}
	})

	parallel.For(0, len(m.workers), 1, func(w int) {
		// TODO(y): Implement Worker.ChangeModel.
		if e := m.workers[w].Call("Worker.UploadDiff", m.aggregators[w%m.cfg.VShards].Addr, nil); e != nil {
			log.Panicf("Worker[%d].UploadDiff(%d) failed: %v", w, w%m.cfg.VShards, e)
		}
	})

	// TODO(y): Not every gibbs iteration need to checkpoint models.
	modelDir := path.Join(tmpDir, "model")
	if e := fs.Mkdir(modelDir); e != nil {
		log.Panicf("Failed to create initialziation model directory %s : %v", modelDir, e)
	}
	parallel.For(0, m.cfg.VShards, 1, func(v int) {
		if e := m.aggregators[v].Call("Aggregator.Save", tmpDir, nil); e != nil {
			log.Panicf("Aggregator %s failed to save model: %v", m.aggregators[v].Addr, e)
		}
	})

	fs.Rename(tmpDir, m.cfg.IterPath(iter))
	log.Printf("Gibbs %d done in %s", iter, time.Since(start))
}

type GibbsArg struct {
	In, Out string
	VShard  int
}

func (m *Master) sampleSegment(segments []string, workers, aggregators []*RPC, iter int) {
	vshards := len(workers)
	for diag := 0; diag <= vshards; diag++ {
		parallel.For(0, vshards, 1, func(v int) {
			s := segments[(v+diag)%vshards]
			var arg = GibbsArg{VShard: v}
			if diag == 0 {
				arg.In = path.Join(m.cfg.IterPath(iter-1), s)
			} else {
				arg.In = path.Join(m.cfg.IterPath(iter), fmt.Sprintf("%s-diag%05d", diag-1))
			}
			if diag == vshards-1 {
				arg.Out = path.Join(m.cfg.IterPath(iter), s)
			} else {
				arg.Out = path.Join(m.cfg.IterPath(iter), fmt.Sprintf("%s-diag%05d", diag))
			}
			if e := workers[v].Call("Worker.SampleVShard", &arg, nil); e != nil {
				log.Panicf("Worker[%d/%s].SampleVShard(%+v): %v", v, workers[v].Addr, arg, e)
			}
		})
	}
}

func (w *Worker) SampleVShard(arg *GibbsArg, _ *int) error {
	in, e := fs.Open(arg.In)
	if e != nil {
		return e
	}
	defer in.Close()

	out, e := fs.Create(arg.Out)
	if e != nil {
		return e
	}
	defer out.Close()

	de := gob.NewDecoder(in)
	en := gob.NewEncoder(out)

	for {
		var d algo.Document
		if e := de.Decode(&d); e != nil {
			if e == io.EOF {
				break
			} else {
				return e
			}
		}
		if e := w.sampler.SampleDocument(&d, arg.VShard); e != nil {
			return e
		}
		if e := en.Encode(&d); e != nil {
			return e
		}
	}

	return nil
}
