package srvs

import (
	"flag"
	"fmt"
	"path"
	"regexp"
	"strconv"
)

type Config struct {
	BaseDir   string
	CorpusDir string
	Vocab     string
	Segmenter string

	Iters   int
	Topics  int
	VShards int
	Groups  int
}

func (cfg *Config) RegisterFlags() {
	flag.StringVar(&cfg.BaseDir, "base", "", "The base directory of a job, as well as job id")
	flag.StringVar(&cfg.CorpusDir, "corpus", "", "The directory in which each file is a training shard")
	flag.StringVar(&cfg.Vocab, "vocab", "", "The token frequency file. Not listed tokens in corpus are ignored.")
	flag.StringVar(&cfg.Segmenter, "segmenter", "", "The segmenter dictionary file.")

	flag.IntVar(&cfg.Iters, "iters", 10, "The number of Gibbs sampling iterations")
	flag.IntVar(&cfg.Topics, "topics", 2, "The number of topics we are going to learn")
	flag.IntVar(&cfg.VShards, "vshards", 2, "The number of VShards of the model")
	flag.IntVar(&cfg.Groups, "groups", 1, "The minimum number of worker groups")
}

// IterDir(0) is the output of initialization iteration. IterDir(1)
// and etc are outputs of Gibbs sampling iterations.  IterDir(i) has
// the form $cfg.BaseDir/iter-0000i.  Please refer to IsIterDir.
func (cfg *Config) IterPath(iter int) string {
	return path.Join(cfg.BaseDir, IterDir(iter))
}

func IterDir(iter int) string {
	return fmt.Sprintf("iter-%05d", iter)
}

func IsIterDir(dir string) bool {
	m, e := regexp.MatchString("^iter-[0-9]+$", dir)
	return e == nil && m
}

func IterFromDir(dir string) (int, error) {
	if !IsIterDir(dir) {
		return -2, fmt.Errorf("%s is not a valid IterDir", dir)
	}
	return strconv.Atoi(dir[5:])
}

func (cfg *Config) VShardPath(iter, vshard int) string {
	return path.Join(cfg.BaseDir, IterDir(iter), VShardName(vshard, cfg.VShards))
}

func VShardName(vshard, vshards int) string {
	return fmt.Sprintf("model/%05d-of-%05d", vshard, vshards)
}

func IsVShardName(name string) bool {
	m, e := regexp.MatchString("^model/[0-9]+-of-[0-9]+$", name)
	return e == nil && m
}
