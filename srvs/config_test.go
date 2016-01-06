package srvs

import (
	"path"
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
