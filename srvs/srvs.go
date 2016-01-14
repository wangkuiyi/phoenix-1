package srvs

import (
	"sync"

	_ "github.com/wangkuiyi/healthz"
	"github.com/wangkuiyi/phoenix/algo"
	"github.com/wangkuiyi/sego"
)

// Worker is an RPC service type.
type Worker struct {
	addr string // Worker network address.
	cfg  Config // Filled by RegisterWorker.

	vocab *algo.Vocab
	vshdr *algo.VSharder
	model *algo.Model

	sgmt *sego.Segmenter
}

// Aggregator is an RPC service type.
type Aggregator struct {
	addr string // Aggregator network address.
	cfg  Config // Filled by RegisterAggregator.

	vocab *algo.Vocab
	vshdr *algo.VSharder
	model *algo.Model
}

// Master is a RPC type.  It waits for workers and aggregators to
// register themselves, and works with workers and aggregators training.
type Master struct {
	cfg               *Config
	workers           []*RPC
	aggregators       []*RPC
	registrationMutex sync.Mutex
	registrationDone  chan bool

	corpusShards []string // Basenames of corpus shard files.
}

func NewMaster(cfg *Config) *Master {
	return &Master{
		cfg:              cfg,
		registrationDone: make(chan bool)}
}
