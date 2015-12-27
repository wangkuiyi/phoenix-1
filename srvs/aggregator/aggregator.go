package aggregator

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"time"

	"github.com/wangkuiyi/healthz"
	"github.com/wangkuiyi/phoenix/srvs"
)

func Run(master string, timeout int) {
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
		if e := srvs.Call(master, "Registry.AddAggregator", w.addr, &w.cfg); e != nil {
			log.Fatalf("Worker %v Cannot register to master: %v", w.addr, e)
		}
	}()

	http.Serve(l, nil)
}

type Aggregator struct {
	addr string
	cfg  *srvs.Config
}

func (a *Aggregator) Aggregate(_ int, _ *int) error {
	return nil
}
