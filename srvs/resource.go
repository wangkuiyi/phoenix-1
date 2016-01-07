package srvs

import (
	"github.com/wangkuiyi/fs"
	"github.com/wangkuiyi/phoenix/algo"
	"github.com/wangkuiyi/sego"
)

func GuaranteeSegmenter(sgmt **sego.Segmenter, dictFile string) error {
	if *sgmt == nil {
		s := new(sego.Segmenter)
		if e := s.LoadDictionary(dictFile); e != nil {
			return e
		}
		*sgmt = s
	}
	return nil
}

func GuaranteeVocabSharder(vocab **algo.Vocab, vshdr **algo.VSharder, tokensFile string, vshards int) error {
	if *vocab == nil || *vshdr == nil {
		tokens, e := fs.Open(tokensFile)
		if e != nil {
			return e
		}
		defer tokens.Close()

		v, s, e := algo.BuildVocabAndVSharder(tokens, vshards, true)
		if e != nil {
			return e
		}
		*vocab = v
		*vshdr = s
	}
	return nil
}
