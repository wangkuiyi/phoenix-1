package master

import (
	"log"

	"github.com/wangkuiyi/fs"
	"github.com/wangkuiyi/phoenix/srvs"
)

type Workflow struct {
	cfg *srvs.Config
}

func NewWorkflow(cfg *srvs.Config) *Workflow {
	return &Workflow{
		cfg: cfg,
	}
}

func (wf *Workflow) Start() {
	if e := fs.Mkdir(wf.cfg.BaseDir); e != nil {
		log.Fatalf("Cannot create job base dir %s: %v", wf.cfg.BaseDir, e)
	}
}
