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
	"github.com/wangkuiyi/phoenix/srvs"
)

func main() {
	basedir := flag.String("basedir", "", "training job basedir")
	iter := flag.Int("iter", 0, "the iteration of whose intermediate corpus are to be dumped")
	flag.Parse()

	cfg := &srvs.Config{BaseDir: *basedir}
	if e := srvs.GuaranteeConfig(cfg); e != nil {
		log.Fatal(e)
	}

	fis, e := fs.ReadDir(cfg.IterPath(*iter))
	if e != nil {
		log.Fatal(e)
	}

	var vocab *algo.Vocab
	var vshdr *algo.VSharder
	if e := srvs.GuaranteeVocabSharder(&vocab, &vshdr, cfg.Vocab, cfg.VShards); e != nil {
		log.Fatal(e)
	}

	for _, fi := range fis {
		f, e := fs.Open(path.Join(cfg.IterPath(*iter), fi.Name()))
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
