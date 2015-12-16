package srvs

import "net/rpc"

type RPC struct {
	*rpc.Client
	Addr string
}
