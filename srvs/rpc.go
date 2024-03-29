package srvs

import (
	"bytes"
	"encoding/gob"
	"net/rpc"
)

type RPC struct {
	*rpc.Client
	Addr string
}

func (c *RPC) GobDecode(buf []byte) error {
	return gob.NewDecoder(bytes.NewBuffer(buf)).Decode(&(c.Addr))
}
func (c *RPC) GobEncode() ([]byte, error) {
	var buf bytes.Buffer
	e := gob.NewEncoder(&buf).Encode(c.Addr)
	return buf.Bytes(), e
}

func Dial(addr string) (*RPC, error) {
	if c, e := rpc.DialHTTP("tcp", addr); e != nil {
		return nil, e
	} else {
		return &RPC{Client: c, Addr: addr}, nil
	}
}

func (c *RPC) Dial() error {
	l, e := rpc.DialHTTP("tcp", c.Addr)
	if e != nil {
		return e
	}
	c.Client = l
	return nil
}

func Call(addr, method string, args, reply interface{}) error {
	c, e := rpc.DialHTTP("tcp", addr)
	if e != nil {
		return e
	}
	defer c.Close()

	return c.Call(method, args, reply)
}
