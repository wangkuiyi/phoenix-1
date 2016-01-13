package srvs

import (
	"crypto/md5"
	"math/rand"
	"unsafe"
)

func Seed(text string) int64 {
	hasher := md5.New()
	hasher.Write([]byte(text))
	sum := hasher.Sum(nil)
	return *(*int64)(unsafe.Pointer(&sum[0]))
}

func NewRand(seed string) *rand.Rand {
	return rand.New(rand.NewSource(Seed(seed)))
}
