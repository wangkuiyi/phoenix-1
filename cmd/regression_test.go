package cmd

import (
	"fmt"
	"os"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/wangkuiyi/fs"
)

func init() {
	buildBinary()
}

func TestRegression(t *testing.T) {
	logDir := "/tmp/phoenix_regression_test"
	os.RemoveAll(logDir)

	base := fmt.Sprintf("/tmp/%s-%d", "phoenix-regression-test", time.Now().UnixNano())

	// NOTE: Here we assuem that :16060 was not used by other programs.
	go runServer(t, "worker", "-master=:16060", "-registration=5", "-log_dir="+logDir)
	time.Sleep(time.Second)
	go runServer(t, "worker", "-master=:16060", "-registration=5", "-log_dir="+logDir)
	time.Sleep(time.Second)
	go runServer(t, "aggregator", "-master=:16060", "-registration=5", "-log_dir="+logDir)
	time.Sleep(time.Second)
	go runServer(t, "master", "-master=:16060", "-base="+base, "-vshards=1", "-groups=2", "-registration=5",
		"-log_dir="+logDir,
		"-segmenter="+path.Join(GoSrc(), "github.com/wangkuiyi/sego/data/dictionary.txt"),
		"-vocab="+path.Join(GoSrc(), "github.com/wangkuiyi/phoenix/algo/testdata/internet-zh.num"),
		"-corpus="+path.Join(GoSrc(), "github.com/wangkuiyi/phoenix/srvs/testdata/corpus"))

	// NOTE: Here we assume that the master will create the base
	// dir once all required servers register themselves.
	time.Sleep(10 * time.Second)
	_, e := fs.Stat(base)
	assert.Nil(t, e)

}
