// Load token frequency list and generates Vocab and VSharder
// file. Example usage:
//
//   $GOPATH/bin/vshard -tf='../testdata/internet-zh.num' -s 100 -vocab /tmp/v -vshdr /tmp/s
//
package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"log"

	"github.com/wangkuiyi/fs"
	"github.com/wangkuiyi/phoenix/algo"
)

func main() {
	freq := flag.String("tf", "", "Token frequency list. Each line consists of frequence and token.")
	shards := flag.Int("s", 10, "Hint of number of vshards. Might be less.")
	delUnbalanced := flag.Bool("d", true, "Delete singular and unbalanced vshards.")
	vocab := flag.String("vocab", "", "Gob encoded Vocab file")
	vshdr := flag.String("vshdr", "", "Gob encoded VSharder file.")
	flag.Parse()

	f, e := fs.Open(*freq)
	if e != nil {
		log.Fatal(e)
	}
	defer f.Close()

	v, e := fs.Create(*vocab)
	if e != nil {
		log.Fatal(e)
	}
	defer v.Close()

	s, e := fs.Create(*vshdr)
	if e != nil {
		log.Fatal(e)
	}
	defer s.Close()

	{
		vocab, vshdr, e := algo.BuildVocabAndVSharder(f, *shards, *delUnbalanced)
		if e != nil {
			log.Fatal(e)
		}

		fmt.Printf("Generated %v shards\n", vshdr.Num())

		if e := gob.NewEncoder(v).Encode(vocab); e != nil {
			log.Fatal(e)
		}

		if e := gob.NewEncoder(s).Encode(vshdr); e != nil {
			log.Fatal(e)
		}
	}
}
