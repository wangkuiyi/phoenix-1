package srvs

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"time"

	"github.com/wangkuiyi/healthz"
)

func RunMaster(addr string, timeout int, cfg *Config) {
	sr := NewRegistry(cfg)
	rpc.Register(sr)
	rpc.HandleHTTP()

	go func() {
		select {
		case <-sr.completion:
			log.Println("Finished server registration. Starting workflow.")
			wf := NewMaster(cfg, sr)
			defer shutdown(wf)
			wf.Start()
		case <-time.After(time.Duration(timeout) * time.Second):
			log.Fatal("Server registration timed out.")
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
// Start and functions called by it cna just panic or log.Panic if
// anything goes wrong.  And, when master restarts, Start uses
// mostRecentCompletedIter to resume training.
func (wf *Master) Start() {
	start := mostRecentCompletedIter(wf.BaseDir)
	wf.Bootstrap(start)

	for i := start; i < wf.Iters; i = mostRecentCompletedIter(wf.BaseDir) {
		if i < 0 {
			wf.Initialize()
		} else {
			wf.Gibbs(i + 1)
		}
	}
}

func RunWorker(master string, timeout int) {
	l, e := net.Listen("tcp", ":0") // OS allocates a free port.
	if e != nil {
		log.Fatalf("Worker cannot listen on: %v", e)
	}

	w := &Worker{addr: l.Addr().String()}
	rpc.Register(w)
	rpc.HandleHTTP()

	go func() {
		if e := healthz.OK(master, time.Duration(timeout)*time.Second); e != nil {
			log.Fatalf("Waiting for master timed out: %v", e)
		}
		if e := Call(master, "Registry.AddWorker", w.addr, &w.cfg); e != nil {
			log.Fatalf("Worker %v Cannot register to master: %v", w.addr, e)
		}
	}()

	http.Serve(l, nil)
}

func RunAggregator(master string, timeout int) {
	l, e := net.Listen("tcp", ":0") // OS allocates a free port.
	if e != nil {
		log.Fatalf("Aggregator cannot listen on: %v", e)
	}

	w := &Aggregator{addr: l.Addr().String()}
	rpc.Register(w)
	rpc.HandleHTTP()

	go func() {
		if e := healthz.OK(master, time.Duration(timeout)*time.Second); e != nil {
			log.Fatalf("Waiting for master timed out: %v", e)
		}
		if e := Call(master, "Registry.AddAggregator", w.addr, &w.cfg); e != nil {
			log.Fatalf("Worker %v Cannot register to master: %v", w.addr, e)
		}
	}()

	http.Serve(l, nil)
}
