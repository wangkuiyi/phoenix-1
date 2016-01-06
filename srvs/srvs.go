package srvs

import (
	_ "github.com/wangkuiyi/healthz"
	"github.com/wangkuiyi/phoenix/algo"
	"github.com/wangkuiyi/sego"
)

// Worker is an RPC service type.
type Worker struct {
	addr string // Worker address. Also worker ID.
	*Config

	vocab *algo.Vocab
	vshdr *algo.VSharder
	sgmt  *sego.Segmenter
}

// Aggregator is an RPC service type.
type Aggregator struct {
	addr string
	*Config

	vocab *algo.Vocab
	vshdr *algo.VSharder
	model *algo.Model
}

// Master is NOT an RPC service type. Instead, Registry is.
type Master struct {
	*Config
	*Registry
}
