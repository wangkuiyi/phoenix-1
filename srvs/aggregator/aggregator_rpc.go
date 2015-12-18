package main

import "github.com/wangkuiyi/phoenix/srvs"

type AggregatorRPC struct {
	addr string
	cfg  *srvs.Config
}

func (a *AggregatorRPC) Aggregate(_ int, _ *int) error {
	return nil
}
