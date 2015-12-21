package master

import (
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"testing"
	"time"

	"github.com/wangkuiyi/parallel"
	"github.com/wangkuiyi/phoenix/srvs"
)

func TestRegistry(t *testing.T) {
	cfg := &srvs.Config{
		VShards: 1,
		Groups:  2}
	sr := NewRegistry(cfg)

	rpc.Register(sr)
	rpc.HandleHTTP()

	if l, e := net.Listen("tcp", ":0"); e != nil {
		t.Error("Master cannot start listening: ", e)
	} else {
		go http.Serve(l, nil)

		go runClients("Worker", 2, l.Addr().String(), t)     // 2 workers
		go runClients("Aggregator", 1, l.Addr().String(), t) // 1 aggregator

		select {
		case <-sr.completion:
		case <-time.After(1 * time.Second):
			t.Errorf("Timeout before registration")
		}
	}
}

func runClients(role string, workers int, master string, t *testing.T) {
	parallel.For(0, workers, 1, func(i int) {
		go func() {
			l, e := net.Listen("tcp", ":0") // OS allocates a free port.
			if e != nil {
				t.Skip("Mocking %s server cannot listen: ", role, e)
			}
			go http.Serve(l, nil)

			if client, e := rpc.DialHTTP("tcp", master); e != nil {
				t.Errorf("%v (%v) cannot dial master (%v): %v", role, l.Addr(), master, e)
			} else {
				var cfg srvs.Config
				e = client.Call(fmt.Sprintf("Registry.Add%v", role), l.Addr().String(), &cfg)
				if e != nil {
					t.Errorf("%v (%v) failed with RPC: %v", role, l.Addr(), e)
				}
			}
		}()
	})
}
