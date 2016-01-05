package srvs

import (
	_ "github.com/wangkuiyi/healthz"
	"github.com/wangkuiyi/phoenix/algo"
	"github.com/wangkuiyi/sego"
)

// Worker is an RPC service type.
type Worker struct {
	addr  string // Worker address. Also worker ID.
	sgmt  *sego.Segmenter
	vocab *algo.Vocab
	vshdr *algo.VSharder
	cfg   *Config
}

// Aggregator is an RPC service type.
type Aggregator struct {
	addr string
	cfg  *Config
}

// Master is NOT an RPC service type. Instead, Registry is.
type Master struct {
	*Config
	*Registry
}
