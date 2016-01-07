package algo

import (
	"bytes"
	"fmt"
	"math/rand"
	"unicode"

	"github.com/wangkuiyi/phoenix/algo/hist"
	"github.com/wangkuiyi/sego"
)

type Document struct {
	VShards   [][]Word
	TopicHist hist.Sparse // TODO(y): To see if we need to improve the data structure.
}

type Word struct {
	Id    int
	Topic int
}

func NewDocument(text string, sgmt *sego.Segmenter, vocab *Vocab, vshdr *VSharder, rng *rand.Rand, topics int, vshards []*Model) *Document {
	d := &Document{
		VShards:   make([][]Word, vshdr.Num()),
		TopicHist: make(hist.Sparse)}

	for _, seg := range sgmt.Segment([]byte(text)) {
		if word := seg.Token().Text(); !allPunctOrSpace(word) {
			if id := vocab.Id(word); id >= 0 {
				shard := vshdr.Shard(id)
				if d.VShards[shard] == nil {
					d.VShards[shard] = make([]Word, 0)
				}
				topic := rng.Intn(topics)
				d.VShards[shard] = append(d.VShards[shard], Word{Id: id, Topic: topic})
				d.TopicHist[topic]++
				if vshards != nil {
					vshards[shard].Inc(id, topic, 1)
				}
			}
		}
	}
	return d
}

func allPunctOrSpace(s string) bool {
	for _, u := range s {
		if !unicode.IsPunct(u) && !unicode.IsSpace(u) {
			return false
		}
	}
	return true
}

func (d *Document) String(vocab *Vocab) string {
	var w bytes.Buffer
	fmt.Fprintf(&w, "{")
	for _, shard := range d.VShards {
		fmt.Fprintf(&w, "[")
		for _, word := range shard {
			fmt.Fprintf(&w, "\"%s\":%d", vocab.Token(word.Id), word.Topic)
		}
		fmt.Fprintf(&w, "]")
	}
	fmt.Fprintf(&w, "%v", d.TopicHist)
	fmt.Fprintf(&w, "}")
	return w.String()
}
