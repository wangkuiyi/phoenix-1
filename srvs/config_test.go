package srvs

import (
	"path"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsIterDir(t *testing.T) {
	assert := assert.New(t)
	cfg := Config{BaseDir: "/tmp/"}
	assert.True(IsIterDir(path.Base(cfg.IterDir(0))))
	assert.True(IsIterDir(path.Base(cfg.IterDir(1))))
	assert.True(IsIterDir(path.Base(cfg.IterDir(10))))
}

func TestIterDirOrder(t *testing.T) {
	assert := assert.New(t)
	cfg := Config{BaseDir: "/tmp/"}

	n := 1001
	s := make([]string, n)
	for i := 0; i < n; i++ {
		s[i] = path.Base(cfg.IterDir(i))
	}
	sort.Strings(s)
	assert.Equal("iter-01000", s[1000])
	assert.Equal("iter-00999", s[999])
	assert.Equal("iter-00001", s[1])
	assert.Equal("iter-00000", s[0])
}
