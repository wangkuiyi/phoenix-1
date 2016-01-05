package master

import (
	"log"
	"os"
	"path"
	"regexp"
	"sort"

	"github.com/wangkuiyi/fs"
	"github.com/wangkuiyi/parallel"
	"github.com/wangkuiyi/phoenix/srvs"
)

type Workflow struct {
	*srvs.Config
	*Registry
}

func NewWorkflow(cfg *srvs.Config, sr *Registry) *Workflow {
	return &Workflow{
		Config:   cfg,
		Registry: sr,
	}
}

// Start is called by master.Run, which recovers panics. Therefore,
// Start and functions called by it cna just panic or log.Panic if
// anything goes wrong.  And, when master restarts, Start uses
// mostRecentCompletedIter to resume training.
func (wf *Workflow) Start() {
	start := mostRecentCompletedIter(wf.BaseDir)
	wf.LoadModel(start)

	for i := start; i < wf.Iters; i = mostRecentCompletedIter(wf.BaseDir) {
		if i < 0 {
			wf.Initialize()
		} else {
			wf.Gibbs(i + 1)
		}
	}
}

// LoadModel makes aggregators loads the existing model if there is
// any, or initialize a new model.
func (wf *Workflow) LoadModel(iter int) {
}

func (wf *Workflow) Initialize() {
	if fis, e := fs.ReadDir(wf.CorpusDir); e != nil {
		log.Panicf("Failed listing corpus shards in %s: %v", wf.CorpusDir, e)
	} else {
		ch := make(chan string)
		parallel.For(0, len(wf.workers), 1, func(i int) {
			for fn := range ch {
				var dumb int
				if e := wf.workers[i].Call("Worker.Initialize", fn, &dumb); e != nil {
					log.Panicf("Worker %s failed initializing shard %s: %v", wf.workers[i].Addr, fn, e)
				}
			}
		})
		for _, fi := range fis {
			if !fi.IsDir() {
				ch <- path.Join(wf.CorpusDir, fi.Name())
			}
		}
		close(ch)
	}
}

func (wf *Workflow) Gibbs(iter int) {
}

// Returns -1 indicating that no iterations is completed, 0 means
// initiialization is completed, 1 means the first Gibbs sampling
// iteration is completed, etc.
func mostRecentCompletedIter(baseDir string) int {
	if base, e := fs.Stat(baseDir); os.IsNotExist(e) {
		if e := fs.Mkdir(baseDir); e != nil {
			log.Panicf("Cannot create job base dir %s: %v", baseDir, e)
		}
		return -1
	} else if !base.IsDir() {
		log.Panicf("%s is specified as the base dir, but it is not a directory", baseDir)
	}

	subdirs, e := fs.ReadDir(baseDir)
	if e != nil {
		log.Panicf("Cannot read base dir %s: %v", baseDir, e)
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
