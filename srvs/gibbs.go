package srvs

import (
	"log"
	"time"
)

func (wf *Master) gibbs(iter int) {
	log.Printf("Gibbs %d ...", iter)
	start := time.Now()

	log.Printf("Gibbs %d done in %s", iter, time.Since(start))
}
