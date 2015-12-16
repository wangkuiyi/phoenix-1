package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"time"

	"github.com/wangkuiyi/healthz"
	"github.com/wangkuiyi/phoenix/srvs"
)

func main() {
	master := flag.String("master", "", "Master address")
	flag.Parse()

	l, e := net.Listen("tcp", ":0") // OS allocates a free port.
	if e != nil {
		log.Fatalf("Mocking aggregator server cannot listen: %v", e)
	}

	w := &WorkerRPC{addr: l.Addr().String()}
	rpc.Register(w)

	go func() {
		if e := healthz.OK(*master, 5*time.Second); e != nil {
			log.Fatalf("Waiting for master timeed out: %v", e)
		}
		if e := srvs.Call(*master, "Registry.AddWorker", w.addr, &w.cfg); e != nil {
			log.Fatalf("Worker %v Cannot register to master: %v", w.addr, e)
		}
	}()

	http.Serve(l, nil)
}
