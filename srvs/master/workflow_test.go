package master

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMostRecentCompletedIter(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(-1, mostRecentCompletedIter("/tmp"))
}
