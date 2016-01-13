package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"io"
	"log"
	"path"

	"github.com/wangkuiyi/fs"
	"github.com/wangkuiyi/phoenix/algo"
	"github.com/wangkuiyi/phoenix/algo/hist"
	"github.com/wangkuiyi/phoenix/srvs"
)

func main() {
	basedir := flag.String("basedir", "", "training job basedir")
	iter := flag.Int("iter", 0, "the iteration of whose intermediate corpus are to be dumped")
	corpus := flag.Bool("corpus", false, "Dump corpus")
	model := flag.Bool("model", true, "Dump model")
	flag.Parse()

	cfg := &srvs.Config{BaseDir: *basedir}
	if e := srvs.GuaranteeConfig(cfg); e != nil {
		log.Fatal(e)
	}

	var vocab *algo.Vocab
	var vshdr *algo.VSharder
	if e := srvs.GuaranteeVocabSharder(&vocab, &vshdr, cfg.Vocab, cfg.VShards); e != nil {
		log.Fatal(e)
	}

	if *corpus {
		dumpCorpus(cfg, *iter, vocab)
	}
	if *model {
		dumpModel(cfg, *iter)
	}
}

func dumpCorpus(cfg *srvs.Config, iter int, vocab *algo.Vocab) {
	fis, e := fs.ReadDir(cfg.IterPath(iter))
	if e != nil {
		log.Fatal(e)
	}

	for _, fi := range fis {
		if fi.IsDir() {
			continue
		}
		f, e := fs.Open(path.Join(cfg.IterPath(iter), fi.Name()))
		if e != nil {
			log.Fatal(e)
		}
		defer f.Close()

		fmt.Println(fi.Name())
		de := gob.NewDecoder(f)
		for {
			d := new(algo.Document)
			if e := de.Decode(d); e != nil {
				if e == io.EOF {
					break
				} else {
					log.Fatal(e)
				}
			} else {
				fmt.Println(d.String(vocab))
			}
		}
	}
}

func dumpModel(cfg *srvs.Config, iter int) {
	for v := 0; v < cfg.VShards; v++ {
		in, e := fs.Open(cfg.VShardPath(iter, v))
		if e != nil {
			log.Fatal(e)
		}
		defer in.Close()

		m := &algo.Model{}
		if e := gob.NewDecoder(in).Decode(m); e != nil {
			log.Fatal(e)
		}

		// TODO(y): Implement the sparse case.
		if m.DenseHists != nil {
			for off, dense := range m.DenseHists {
				if len(dense) > 0 {
					o := hist.OrderedFromDense(dense)
					if len(o.Topics) > 0 {
						fmt.Printf("%s : %s\n", m.Vocab.Token(off+m.VShdr.Begin(m.VShard)), o)
					}
				}
			}
		}
	}
}
