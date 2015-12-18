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
	timeout := flag.Int("registration", 5, "Registeration timeout in seconds")
	flag.Parse()

	l, e := net.Listen("tcp", ":0") // OS allocates a free port.
	if e != nil {
		log.Fatalf("Worker cannot listen on: %v", e)
	}

	w := &WorkerRPC{addr: l.Addr().String()}
	rpc.Register(w)
	rpc.HandleHTTP()

	go func() {
		if e := healthz.OK(*master, time.Duration(*timeout)*time.Second); e != nil {
			log.Fatalf("Waiting for master timed out: %v", e)
		}
		if e := srvs.Call(*master, "Registry.AddWorker", w.addr, &w.cfg); e != nil {
			log.Fatalf("Worker %v Cannot register to master: %v", w.addr, e)
		}
	}()

	http.Serve(l, nil)
}
