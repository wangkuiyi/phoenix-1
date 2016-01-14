package srvs

import (
	"log"
	"time"

	"github.com/wangkuiyi/parallel"
)

func (m *Master) gibbs(iter int) {
	log.Printf("Gibbs %d ...", iter)
	start := time.Now()

	ch := make(chan []string)
	go func() {
		for s := 0; s < len(m.corpusShards); s += m.cfg.VShards {
			e := s + m.cfg.VShards
			if e > len(m.corpusShards) {
				e = len(m.corpusShards)
			}
			ch <- m.corpusShards[s:e]
		}
	}()
	parallel.For(0, len(m.workers), m.cfg.VShards, func(w int) {
		for seg := range ch {
			sampleSegment(seg, m.workers[w:w+m.cfg.VShards], m.aggregators)
		}
	})

	log.Printf("Gibbs %d done in %s", iter, time.Since(start))
}
