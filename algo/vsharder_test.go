package algo

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wangkuiyi/fs"
)

func BenchmarkBuildVocabAndVSharder(t *testing.B) {
	f, e := fs.Open("testdata/internet-zh.num")
	if e != nil {
		t.Skip(e)
	}
	defer f.Close()

	_, _, e = BuildVocabAndVSharder(f, 10, true)
	if e != nil {
		t.Skip(e)
	}
}

func TestBuildVocabAndVSharder(t *testing.T) {
	content := `0.1 我
1.1 的
0.2 ，
0.1 们
`
	v, s, e := BuildVocabAndVSharder(strings.NewReader(content), 3, true)
	if e != nil {
		t.Skip(e)
	}

	assert := assert.New(t)
	assert.Equal(s.Num(), 2)
	assert.Equal(s.Size(0), 2)
	assert.Equal(s.Size(1), 1)
	assert.Equal(s.Begin(0), 0)
	assert.Equal(s.Begin(1), 2)

	assert.Equal("我", v.Token(0))
	assert.Equal("们", v.Token(1))
	assert.Equal("，", v.Token(2))

	assert.Equal(0, v.Id("我"))
	assert.Equal(-1, v.Id("的"))
	assert.Equal(2, v.Id("，"))
	assert.Equal(1, v.Id("们"))

	assert.Equal(0, s.Shard(v.Id("我")))
	assert.Equal(-1, s.Shard(v.Id("的")))
	assert.Equal(1, s.Shard(v.Id("，")))
	assert.Equal(0, s.Shard(v.Id("们")))
}
