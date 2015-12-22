package worker

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"time"

	"github.com/wangkuiyi/fs"
	"github.com/wangkuiyi/healthz"
	"github.com/wangkuiyi/phoenix/algo"
	"github.com/wangkuiyi/phoenix/srvs"
	"github.com/wangkuiyi/sego"
)

func Run(master string, timeout int) {
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
		if e := srvs.Call(master, "Registry.AddWorker", w.addr, &w.cfg); e != nil {
			log.Fatalf("Worker %v Cannot register to master: %v", w.addr, e)
		}
	}()

	http.Serve(l, nil)
}

// Worker is the RPC service implementation.
type Worker struct {
	addr  string // Worker address. Also worker ID.
	sgmt  *sego.Segmenter
	vocab *algo.Vocab
	vshdr *algo.VSharder
	cfg   *srvs.Config
}

func (w *Worker) Initialize(shard_filename string, _ *int) error {
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
