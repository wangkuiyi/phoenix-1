package srvs

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/wangkuiyi/fs"
)

func init() {
	buildBinary()
}

func buildBinary() {
	b, e := exec.Command("go", "install", "github.com/wangkuiyi/phoenix/cmd/ssh").CombinedOutput()
	if e != nil {
		log.Panicf("Failed building: %v: %s", e, b)
	}
}

func TestRegression(t *testing.T) {
	base := fmt.Sprintf("/tmp/%s-%d", "phoenix-regression-test", time.Now().UnixNano())

	// NOTE: Here we assuem that :16060 was not used by other programs.
	go runServer(t, "worker", "-master=:16060", "-registration=5", "-logPrefix=/tmp/regression_test")
	time.Sleep(time.Second)
	go runServer(t, "worker", "-master=:16060", "-registration=5", "-logPrefix=/tmp/regression_test")
	time.Sleep(time.Second)
	go runServer(t, "aggregator", "-master=:16060", "-registration=5", "-logPrefix=/tmp/regression_test")
	time.Sleep(time.Second)
	go runServer(t, "master", "-master=:16060", "-base="+base, "-vshards=1", "-groups=2", "-registration=5",
		"-logPrefix=/tmp/regression_test",
		"-segmenter="+path.Join(GoSrc(), "github.com/wangkuiyi/sego/data/dictionary.txt"),
		"-vocab="+path.Join(GoSrc(), "github.com/wangkuiyi/phoenix/algo/testdata/internet-zh.num"),
		"-corpus="+path.Join(GoSrc(), "github.com/wangkuiyi/phoenix/srvs/testdata/corpus"))

	// NOTE: Here we assume that the master will create the base
	// dir once all required servers register themselves.
	time.Sleep(10 * time.Second)
	_, e := fs.Stat(base)
	assert.Nil(t, e)

}

func runServer(t *testing.T, role string, args ...string) {
	p := path.Join(GoPath(), "bin", "ssh")
	args = append(args, "-role="+role)
	fmt.Printf("Starting %v %v\n", p, strings.Join(args, " "))
	// NOTE: the sub-process created by Cmd.CombinedOutput will be
	// killed after this test process completes.
	b, e := exec.Command(p, args...).CombinedOutput()
	if e != nil {
		t.Skipf("Failed runing %s: %v: %s", role, e, b)
	}
}

func GoPath() string {
	return strings.Split(os.Getenv("GOPATH"), ":")[0]
}

func GoSrc() string {
	return path.Join(GoPath(), "src")
}
