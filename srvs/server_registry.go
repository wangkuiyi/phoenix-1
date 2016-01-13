package srvs

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/wangkuiyi/healthz"
)

func (sr *Master) done() bool {
	return len(sr.aggregators) == sr.cfg.VShards && len(sr.workers) >= sr.cfg.VShards*sr.cfg.Groups
}

// NOTE: RegisterWorker doesn't put a cap on the number of workers. So
// there might be workers register themselves after training starts.
func (sr *Master) RegisterWorker(addr string, cfg *Config) error {
	sr.registrationMutex.Lock()
	defer sr.registrationMutex.Unlock()

	if c, e := Dial(addr); e == nil {
		log.Printf("Established connection to registered worker %s.", addr)
		sr.workers = append(sr.workers, c)
		*cfg = *sr.cfg
		if sr.done() {
			sr.registrationDone <- true
		}
		return nil
	} else {
		e = fmt.Errorf("Master cannot connect to registering worker %v: %v", addr, e)
		log.Print(e)
		return e
	}
}

func (sr *Master) RegisterAggregator(addr string, cfg *Config) error {
	sr.registrationMutex.Lock()
	defer sr.registrationMutex.Unlock()

	if len(sr.aggregators) >= sr.cfg.VShards {
		return errors.New("No more aggregators required")
	}

	if c, e := Dial(addr); e == nil {
		log.Printf("Established connection to registered aggregator %s.", addr)
		sr.aggregators = append(sr.aggregators, c)
		*cfg = *sr.cfg
		if sr.done() {
			sr.registrationDone <- true
		}
		return nil
	} else {
		e = fmt.Errorf("Master cannot connect to registering aggregator %s: %s", addr, e)
		log.Print(e)
		return e
	}
}

func RegisterWorker(master, worker string, cfg *Config, timeoutSeconds int) {
	if e := healthz.OK(master, time.Duration(timeoutSeconds)*time.Second); e != nil {
		log.Panicf("Waiting for master timed out: %v", e)
	}
	if e := Call(master, "Master.RegisterWorker", worker, cfg); e != nil {
		log.Panicf("Worker %v Cannot register to master: %v", worker, e)
	}
}

func RegisterAggregator(master, aggregator string, cfg *Config, timeoutSeconds int) {
	if e := healthz.OK(master, time.Duration(timeoutSeconds)*time.Second); e != nil {
		log.Panicf("Waiting for master timed out: %v", e)
	}
	if e := Call(master, "Master.RegisterAggregator", aggregator, cfg); e != nil {
		log.Panicf("Worker %v Cannot register to master: %v", aggregator, e)
	}
}
