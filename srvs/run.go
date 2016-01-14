package srvs

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"time"
)

func RunMaster(addr string, timeout int, cfg *Config) {
	m := NewMaster(cfg)
	rpc.Register(m)
	rpc.HandleHTTP()

	go func() {
		select {
		case <-m.registrationDone:
			log.Printf("Finished server registration. Starting training: %v", *cfg)
			defer shutdown(m)
			m.start()
		case <-time.After(time.Duration(timeout) * time.Second):
			log.Panic("Server registration timed out.")
		}
	}()

	if e := http.ListenAndServe(addr, nil); e != nil {
		log.Panic(e)
	}
}

func shutdown(wf *Master) {
	// TODO(y): shutdown all workers and aggregators when master panics.
}

// Start is called by RunMaster, which recovers panics. Therefore,
// Start and functions called by it can just panic or log.Panic if
// anything goes wrong.  And, when master restarts, Start uses
// mostRecentCompletedIter to resume training.
func (m *Master) start() {
	start := mostRecentCompletedIter(m.cfg.BaseDir)
	m.bootstrap(start)

	for i := start; i < m.cfg.Iters; i = mostRecentCompletedIter(m.cfg.BaseDir) {
		log.Println("Iteration ", i)
		if i < 0 {
			m.initialize()
		} else {
			m.gibbs(i + 1)
		}
	}
}

func RunWorker(master string, timeoutSeconds int) {
	l, e := net.Listen("tcp", ":0") // OS allocates a free port.
	if e != nil {
		log.Panicf("Worker cannot listen on: %v", e)
	}
	log.Printf("Worker listening on %s", l.Addr())

	w := &Worker{addr: l.Addr().String()}
	rpc.Register(w)
	rpc.HandleHTTP()

	go func() {
		RegisterWorker(master, w.addr, &w.cfg, timeoutSeconds)
	}()

	http.Serve(l, nil)
}

func RunAggregator(master string, timeoutSeconds int) {
	l, e := net.Listen("tcp", ":0") // OS allocates a free port.
	if e != nil {
		log.Panicf("Aggregator cannot listen on: %v", e)
	}
	log.Printf("Aggregator listening on %s", l.Addr())

	a := &Aggregator{addr: l.Addr().String()}
	rpc.Register(a)
	rpc.HandleHTTP()

	go func() {
		RegisterAggregator(master, a.addr, &a.cfg, timeoutSeconds)
	}()

	http.Serve(l, nil)
}
