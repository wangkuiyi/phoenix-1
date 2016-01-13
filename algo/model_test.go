package algo

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewModel(t *testing.T) {
	assert := assert.New(t)
	content := `0.1 我
1.1 的
0.2 ，
0.1 们
`
	// Build a 2-vshard VSharder. For details please refer to vsharder_test.go.
	v, s, _ := BuildVocabAndVSharder(strings.NewReader(content), 3, true)

	m0 := NewModel(true, v, s, 0, 2)
	m1 := NewModel(false, v, s, 1, 2)

	assert.Equal(2, m0.Topics())
	assert.Equal(2, m1.Topics())

	m0.Inc(v.Id("我"), 0, 1)
	m1.Inc(v.Id("我"), 0, 1)
	m0.Inc(v.Id("的"), 1, 1)
	m1.Inc(v.Id("的"), 1, 1)
	m0.Inc(v.Id("，"), 0, 1)
	m1.Inc(v.Id("，"), 0, 1)
	m0.Inc(v.Id("们"), 1, 1)
	m1.Inc(v.Id("们"), 1, 1)
	m0.Inc(v.Id("a"), 0, 1)
	m1.Inc(v.Id("a"), 0, 1)

	assert.Equal(int32(1), m0.Dense(v.Id("我"))[0])
	assert.Equal(int32(0), m0.Dense(v.Id("我"))[1])
	assert.Equal(int32(0), m0.Dense(v.Id("们"))[0])
	assert.Equal(int32(1), m0.Dense(v.Id("们"))[1])
	assert.Equal(int32(1), m1.Sparse(v.Id("，"))[0])
	assert.Equal(int32(0), m1.Sparse(v.Id("，"))[1])
}
