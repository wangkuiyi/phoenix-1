package srvs

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

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

// GuaranteeConfig writes cfg into cfg.BaseDir/config if the file
// doesn't exist, or load from that file and overwrites cfg.  Anyway,
// GuaranteeConfig requires that cfg.BaseDir is set.
func GuaranteeConfig(cfg *Config) error {
	file := path.Join(cfg.BaseDir, "config")
	if _, e := fs.Stat(file); os.IsNotExist(e) {
		c, e := fs.Create(file)
		if e != nil {
			return fmt.Errorf("Failed to create config file %s: %v", file, e)
		}
		if e := json.NewEncoder(c).Encode(cfg); e != nil {
			return fmt.Errorf("Failed to encode config file %s: %v", file, e)
		}
		c.Close()
	} else {
		c, e := fs.Open(file)
		if e != nil {
			return fmt.Errorf("Failed to open config file %s: %v", file, e)
		}
		if e := json.NewDecoder(c).Decode(cfg); e != nil {
			return fmt.Errorf("Failed to decode config file %s: %v", file, e)
		}
		c.Close()
	}
	return nil
}
