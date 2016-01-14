package srvs

import (
	"log"
	"os"
	"sort"

	"github.com/wangkuiyi/fs"
)

// Returns -1 indicating that no iterations is completed, 0 means
// initiialization is completed, 1 means the first Gibbs sampling
// iteration is completed, etc.  For the form of IterDir, please refer
// to Config.IterDir and IsIterDir.
//
// Note that mostRecentCompletedIter checks the completion of an
// iteration by checking the existence of its IterDir.  So we should
// create a temporary directory and rename it into the form of an
// IterDir.
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
	for _, sd := range subdirs {
		if i, e := IterFromDir(sd.Name()); e == nil {
			// TODO(y): Not every gibbs iteration need to checkpoint models. Should return the most recent iteration with /model/ subdir.
			return i
		}
	}
	return -1
}

type byName []os.FileInfo

func (f byName) Len() int           { return len(f) }
func (f byName) Less(i, j int) bool { return f[i].Name() > f[j].Name() }
func (f byName) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }
