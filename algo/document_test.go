package algo

import (
	"math/rand"
	"os"
	"path"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wangkuiyi/sego"
)

var (
	sgmt sego.Segmenter
)

func init() {
	sgmt.LoadDictionary(path.Join(os.Getenv("GOPATH"), "src/github.com/wangkuiyi/sego/data/dictionary.txt"))
}

// This test shows that sego.Segmenter inserts space at where we
// manually add space.  If we add consecutive spaces, it outputs as
// many spaces we insertd.  \t is also considered a space.  It also
// shows that capitalized Latin letters are lowered.
func TestSegoSegmenter(t *testing.T) {
	segment := func(s string) []string {
		r := make([]string, 0)
		for _, s := range sgmt.Segment([]byte(s)) {
			r = append(r, s.Token().Text())
		}
		return r
	}

	inputs := []string{
		"今天玩QQ poker游戏high吗？",
		"今天玩QQ poker游 戏high吗？",
		"今 天玩QQ poker游  戏high吗？",
		"今 天玩QQ poker 游 \t 戏high吗？"}
	outputs := [][]string{
		{"今天", "玩", "qq", " ", "poker", "游戏", "high", "吗", "？"},
		{"今天", "玩", "qq", " ", "poker", "游", " ", "戏", "high", "吗", "？"},
		{"今", " ", "天", "玩", "qq", " ", "poker", "游", " ", " ", "戏", "high", "吗", "？"},
		{"今", " ", "天", "玩", "qq", " ", "poker", " ", "游", " ", "\t", " ", "戏", "high", "吗", "？"}}
	for i := range inputs {
		if r := segment(inputs[i]); !reflect.DeepEqual(outputs[i], r) {
			t.Errorf("Case %v: expecting %v, when input is %v, but got %v", i, outputs[i], inputs[i], r)
		}
	}
}

// This test shows that the unicode package considers most of Latin
// and Chinese punctuators as punctuators, but not including = and +.
func TestAllPuctOrSpace(t *testing.T) {
	assert := assert.New(t)

	assert.True(allPunctOrSpace("\t  ?？//()（）[]【】{},，.。、——"))
	assert.True(allPunctOrSpace("-"))

	assert.False(allPunctOrSpace("="))
	assert.False(allPunctOrSpace("+"))
}

func TestNewDocument(t *testing.T) {
	inputs := []string{
		"今天玩QQ poker游戏high吗？",
		"今天玩QQ poker游 戏high吗？",
		"今 天玩QQ poker游  戏high吗？",
		"今 天玩QQ poker 游 \t 戏high吗？"}
	outputs := []string{
		`{["今天":0]["玩":0"qq":0"poker":0"游戏":0]["high":0"吗":0]map[0:7]}`,
		`{["今天":0]["玩":0"qq":0"poker":0"游":0"戏":0]["high":0"吗":0]map[0:8]}`,
		`{["今":0"天":0]["玩":0"qq":0"poker":0"游":0"戏":0]["high":0"吗":0]map[0:9]}`,
		`{["今":0"天":0]["玩":0"qq":0"poker":0"游":0"戏":0]["high":0"吗":0]map[0:9]}`}
	vocab := &Vocab{
		Ids:    map[string]int{"今": 0, "天": 1, "今天": 2, "玩": 3, "qq": 4, "poker": 5, "游": 6, "戏": 7, "游戏": 8, "high": 9, "吗": 10, "？": 11},
		Tokens: []string{"今", "天", "今天", "玩", "qq", "poker", "游", "戏", "游戏", "high", "吗", "？"}}
	vshdr := &VSharder{
		Breakpoints: []int{3, 9, 12}}
	rng := rand.New(rand.NewSource(1))
	topics := 1
	for i := range inputs {
		if r := NewDocument(inputs[i], &sgmt, vocab, vshdr, rng, topics, nil).String(vocab); r != outputs[i] {
			t.Errorf("Expecting %v, got %v", outputs[i], r)
		}
	}
}
