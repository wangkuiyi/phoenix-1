package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/rpc"
	"time"

	"github.com/wangkuiyi/phoenix/srvs"
)

func main() {
	addr := flag.String("addr", "", "Master listening address")
	var cfg srvs.Config
	cfg.RegisterFlags()
	flag.Parse()

	sr := NewServerRegistry(cfg.VShards, cfg.MinGroups)

	rpc.Register(sr)
	rpc.HandleHTTP()

	go waitForServerRegistration(sr)
	if e := http.ListenAndServe(*addr, nil); e != nil {
		log.Panic(e)
	}
}

func waitForServerRegistration(sr *ServerRegistry) {
	select {
	case <-sr.completion:
		fmt.Println("Start working ...")
	case <-time.After(5 * time.Second): //TODO(y): make this a config flag
		fmt.Println("Timeout ...")
	}
}
