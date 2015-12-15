package main

import "github.com/wangkuiyi/phoenix/srvs"

type MasterRPC struct {
}

func (m *MasterRPC) Initialize(cfg *srvs.Config, shards *int) error {
	// fs, e := fs.ReadDir(cfg.CorpusDir)
	// if e != nil {
	// 	return fmt.Errorf("master.Initialize(%v): %v", cfg.CorpusDir, e)
	// }

	// ch := make(chan string, 0)
	return nil
}
