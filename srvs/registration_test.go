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
	buildServer("master")
	buildServer("worker")
	buildServer("aggregator")
}

// server could be "master", "worker", or "aggregator".
func buildServer(server string) {
	p := path.Join("github.com/wangkuiyi/phoenix/srvs", server)
	b, e := exec.Command("go", "install", p).CombinedOutput()
	if e != nil {
		log.Panicf("Failed building %s: %v: %s", server, e, b)
	}
}

func TestRegistration(t *testing.T) {
	base := fmt.Sprintf("/tmp/%d-%d", os.Getpid(), time.Now().UnixNano())

	// NOTE: Here we assuem that :16060 was not used by other programs.
	go runServer(t, "worker", "-master=:16060", "-registration=5")
	time.Sleep(time.Second)
	go runServer(t, "worker", "-master=:16060", "-registration=5")
	time.Sleep(time.Second)
	go runServer(t, "aggregator", "-master=:16060", "-registration=5")
	time.Sleep(time.Second)
	go runServer(t, "master", "-addr=:16060", "-base="+base, "-vshards=1", "-groups=2", "-registration=5")

	// NOTE: Here we assume that the master will create the base
	// dir once all required servers register themselves.
	time.Sleep(2 * time.Second)
	_, e := fs.Stat(base)
	assert.Nil(t, e)

}

func runServer(t *testing.T, server string, args ...string) {
	p := path.Join(os.Getenv("GOPATH"), "bin", server)
	fmt.Printf("Starting %v %v\n", p, strings.Join(args, " "))
	// NOTE: the sub-process created by Cmd.CombinedOutput will be kill after this test process completes.
	b, e := exec.Command(p, args...).CombinedOutput()
	if e != nil {
		t.Skipf("Failed runing %s: %v: %s", server, e, b)
	}
}
