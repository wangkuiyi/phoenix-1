package hist

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrdered(t *testing.T) {
	assert.Equal(t, "1:3 3:2 0:1 ", OrderedFromDense([]int32{1, 3, 0, 2}).String())
	assert.Equal(t, "0:11 1:10 2:9 3:8 4:7 5:6 6:5 7:4 8:3 9:2 ", OrderedFromDense([]int32{11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1}).String())
}
