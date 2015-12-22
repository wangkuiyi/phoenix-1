package main

import (
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/wangkuiyi/phoenix/srvs"
)

// Registry is a RPC type.  It allows and expects vshards aggregators
// and at least group*vshards workers to register themselves.
type Registry struct {
	vshards     int
	minGroups   int
	cfg         *srvs.Config
	workers     []*srvs.RPC
	aggregators []*srvs.RPC
	mutex       sync.Mutex
	completion  chan bool
}

// Creates a Registry RPC service, which will trigger channel
// completion after vshards aggregators and at least minGroups*vshards
// workers registered.
func NewRegistry(cfg *srvs.Config) *Registry {
	return &Registry{
		vshards:    cfg.VShards,
		minGroups:  cfg.MinGroups,
		cfg:        cfg,
		completion: make(chan bool)}
}

func (sr *Registry) completed() bool {
	return len(sr.aggregators) == sr.vshards && len(sr.workers) >= sr.vshards*sr.minGroups
}

// NOTE: AddWorker doesn't put a cap on the number of workers. So
// there might be workers register themselves after the workflow
// starts.
func (sr *Registry) AddWorker(addr string, cfg *srvs.Config) error {
	sr.mutex.Lock()
	defer sr.mutex.Unlock()

	if c, e := srvs.Dial(addr); e == nil {
		log.Printf("Established connection to registered worker %s.", addr)
		sr.workers = append(sr.workers, c)
		*cfg = *sr.cfg
		if sr.completed() {
			sr.completion <- true
		}
		return nil
	} else {
		e = fmt.Errorf("Master cannot connect to registering worker %v: %v", addr, e)
		log.Print(e)
		return e
	}
}

func (sr *Registry) AddAggregator(addr string, cfg *srvs.Config) error {
	sr.mutex.Lock()
	defer sr.mutex.Unlock()

	if len(sr.aggregators) >= sr.vshards {
		return errors.New("No more aggregators required")
	}

	if c, e := srvs.Dial(addr); e == nil {
		log.Printf("Established connection to registered aggregator %s.", addr)
		sr.aggregators = append(sr.aggregators, c)
		*cfg = *sr.cfg
		if sr.completed() {
			sr.completion <- true
		}
		return nil
	} else {
		e = fmt.Errorf("Master cannot connect to registering aggregator %s: %s", addr, e)
		log.Print(e)
		return e
	}
}
