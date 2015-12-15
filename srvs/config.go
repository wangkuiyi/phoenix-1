package srvs

import (
	"fmt"
	"path"
)

type Config struct {
	BaseDir   string // This also identifies a job.
	Topics    int    // Number of topics we are learning.
	CorpusDir string // Where corpus shards reside.
}

func (cfg *Config) IterDir(iter int) string {
	return path.Join(cfg.BaseDir, fmt.Sprintf("%05d", iter))
}
