package master

import (
	"log"
	"net/http"
	"net/rpc"
	"time"

	_ "github.com/wangkuiyi/healthz"
	"github.com/wangkuiyi/phoenix/srvs"
)

func Run(addr string, timeout int, cfg *srvs.Config) {
	sr := NewRegistry(cfg)
	rpc.Register(sr)
	rpc.HandleHTTP()

	go func() {
		select {
		case <-sr.completion:
			log.Println("Finished server registration. Starting workflow.")
			NewWorkflow(cfg).Start()
		case <-time.After(time.Duration(timeout) * time.Second):
			log.Fatal("Server registration timed out.")
		}
	}()

	if e := http.ListenAndServe(addr, nil); e != nil {
		log.Panic(e)
	}
}
