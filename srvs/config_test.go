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
	assert.True(IsIterDir(path.Base(cfg.IterPath(0))))
	assert.True(IsIterDir(path.Base(cfg.IterPath(1))))
	assert.True(IsIterDir(path.Base(cfg.IterPath(10))))

	assert.False(IsIterDir(path.Base(cfg.IterPath(-1))))
}

func TestIterDirOrder(t *testing.T) {
	assert := assert.New(t)
	cfg := Config{BaseDir: "/tmp/"}

	n := 1001
	s := make([]string, n)
	for i := 0; i < n; i++ {
		s[i] = path.Base(cfg.IterPath(i))
	}
	sort.Strings(s)
	assert.Equal("iter-01000", s[1000])
	assert.Equal("iter-00999", s[999])
	assert.Equal("iter-00001", s[1])
	assert.Equal("iter-00000", s[0])
}

func TestIterFromDir(t *testing.T) {
	assert := assert.New(t)
	i, _ := IterFromDir("iter-01000")
	assert.Equal(1000, i)
	i, _ = IterFromDir("iter-00000")
	assert.Equal(0, i)
}

func TestVShardName(t *testing.T) {
	assert := assert.New(t)
	assert.Equal("model/00000-of-00005", VShardName(0, 5))

	assert.True(IsVShardName(VShardName(0, 5)))
	assert.False(IsVShardName("odel/00000-of-00005"))
}
