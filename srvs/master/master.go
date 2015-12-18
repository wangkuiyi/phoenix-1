package main

import (
	"flag"
	"log"
	"net/http"
	"net/rpc"
	"time"

	_ "github.com/wangkuiyi/healthz"
	"github.com/wangkuiyi/phoenix/srvs"
)

func main() {
	addr := flag.String("addr", "", "Master listening address")
	timeout := flag.Int("registration", 5, "Registeration timeout in seconds")
	var cfg srvs.Config
	cfg.RegisterFlags()
	flag.Parse()

	sr := NewRegistry(&cfg)
	wf := NewWorkflow(&cfg)

	rpc.Register(sr)
	rpc.HandleHTTP()

	go func() {
		select {
		case <-sr.completion:
			wf.Start()
		case <-time.After(time.Duration(*timeout) * time.Second): //TODO(y): make this a config flag
			log.Fatal("Server registration timed out.")
		}
	}()

	if e := http.ListenAndServe(*addr, nil); e != nil {
		log.Panic(e)
	}
}
