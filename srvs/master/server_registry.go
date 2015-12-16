package main

import (
	"errors"
	"fmt"
	"log"
	"net/rpc"
	"sync"

	"github.com/wangkuiyi/phoenix/srvs"
)

type ServerRegistry struct {
	vshards     int
	minGroups   int
	workers     []*srvs.RPC
	aggregators []*srvs.RPC
	mutex       sync.Mutex
	completion  chan bool
}

// Creates a ServerRegistry RPC service, which will trigger channel
// completion after vshards aggregators and at least minGroups*vshards
// workers registered.
func NewServerRegistry(vshards, minGroups int) *ServerRegistry {
	return &ServerRegistry{
		vshards:    vshards,
		minGroups:  minGroups,
		completion: make(chan bool)}
}

func (sr *ServerRegistry) completed() bool {
	return len(sr.aggregators) == sr.vshards && len(sr.workers) >= sr.vshards*sr.minGroups
}

func (sr *ServerRegistry) RegisterWorker(addr string, connected *bool) error {
	sr.mutex.Lock()
	defer sr.mutex.Unlock()
	*connected = false

	if len(sr.workers) >= sr.vshards*sr.minGroups {
		return errors.New("No more workers required")
	}

	if client, e := rpc.DialHTTP("tcp", addr); e == nil {
		log.Printf("Established connection to registered worker %s.", addr)
		sr.workers = append(sr.workers, &srvs.RPC{Client: client, Addr: addr})
		*connected = true
		if sr.completed() {
			sr.completion <- true
		}
		return nil
	} else {
		e = fmt.Errorf("Cannot connect to registered worker %v: %v", addr, e)
		log.Print(e)
		return e
	}
}

func (sr *ServerRegistry) RegisterAggregator(addr string, connected *bool) error {
	sr.mutex.Lock()
	defer sr.mutex.Unlock()
	*connected = false

	if len(sr.aggregators) >= sr.vshards {
		return errors.New("No more aggregators required")
	}

	if client, e := rpc.DialHTTP("tcp", addr); e == nil {
		log.Printf("Established connection to registered aggregator %s.", addr)
		sr.aggregators = append(sr.aggregators, &srvs.RPC{Client: client, Addr: addr})
		*connected = true
		if sr.completed() {
			sr.completion <- true
		}
		return nil
	} else {
		e = fmt.Errorf("Cannot connect to aggregator %s: %s", addr, e)
		log.Print(e)
		return e
	}
}
