package srvs

import (
	"crypto/md5"
	"math/rand"
	"net/rpc"
	"unsafe"
)

type RPC struct {
	*rpc.Client
	Addr string
}

func Dial(addr string) (*RPC, error) {
	if c, e := rpc.DialHTTP("tcp", addr); e != nil {
		return nil, e
	} else {
		return &RPC{Client: c, Addr: addr}, nil
	}
}

func Call(addr, method string, args, reply interface{}) error {
	c, e := rpc.DialHTTP("tcp", addr)
	if e != nil {
		return e
	}
	defer c.Close()

	return c.Call(method, args, reply)
}

func Seed(text string) int64 {
	hasher := md5.New()
	hasher.Write([]byte("Hello"))
	sum := hasher.Sum(nil)
	return *(*int64)(unsafe.Pointer(&sum[0]))
}

func NewRand(seed string) *rand.Rand {
	return rand.New(rand.NewSource(Seed(seed)))
}
