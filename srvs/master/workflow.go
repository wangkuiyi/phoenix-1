package master

import (
	"log"
	"os"
	"regexp"
	"sort"

	"github.com/wangkuiyi/fs"
	"github.com/wangkuiyi/phoenix/srvs"
)

type Workflow struct {
	*srvs.Config
	sr *Registry
}

func NewWorkflow(cfg *srvs.Config, sr *Registry) *Workflow {
	return &Workflow{
		Config: cfg,
		sr:     sr,
	}
}

func (wf *Workflow) Start() {
	for i := wf.mostRecentCompletedIter(); i < wf.Iters; i = wf.mostRecentCompletedIter() {
		if i < 0 {
			wf.Initialize()
		} else {
			wf.Gibbs(i + 1)
		}
	}
}

// Returns -1 indicating that no iterations is completed, 0 means
// initiialization is completed, 1 means the first Gibbs sampling
// iteration is completed, etc.
func (wf *Workflow) mostRecentCompletedIter() int {
	if base, e := fs.Stat(wf.BaseDir); os.IsNotExist(e) {
		if e := fs.Mkdir(wf.BaseDir); e != nil {
			log.Panicf("Cannot create job base dir %s: %v", wf.BaseDir, e)
		}
		return -1
	} else if !base.IsDir() {
		log.Panicf("%s is specified as the base dir, but it is not a directory", wf.BaseDir)
	}

	subdirs, e := fs.ReadDir(wf.BaseDir)
	if e != nil {
		log.Panicf("Cannot read base dir %s: %v", wf.BaseDir, e)
	}
	sort.Sort(byName(subdirs)) // sort by descending alphabetic order of names
	for i, sd := range subdirs {
		if m, e := regexp.MatchString("^foo[0-9]+$", sd.Name()); e != nil && m {
			return i
		}
	}
	return -1
}

type byName []os.FileInfo

func (f byName) Len() int           { return len(f) }
func (f byName) Less(i, j int) bool { return f[i].Name() > f[j].Name() }
func (f byName) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }

func (wf *Workflow) Initialize() {
}

func (wf *Workflow) Gibbs(iter int) {
}
