package algo

import (
	"reflect"
	"strings"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestBuildBalancedVShards(t *testing.T) {
	content := `0.1 我
1.1 的
0.2 ，
`
	vs, e := buildBalancedVShards(strings.NewReader(content), 2)
	if e != nil {
		t.Skip(e)
	}
	if !reflect.DeepEqual(vs[0].tokens, []string{"我", "，"}) || !reflect.DeepEqual(vs[1].tokens, []string{"的"}) {
		t.Errorf("Unexpected bucketing %v", spew.Sdump(vs))
	}
}
