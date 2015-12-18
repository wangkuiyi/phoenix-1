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
	go runServer(t, "worker", "-master=:6060", "-registration=5")
	time.Sleep(time.Second)
	go runServer(t, "worker", "-master=:6060", "-registration=5")
	time.Sleep(time.Second)
	go runServer(t, "aggregator", "-master=:6060", "-registration=5")
	time.Sleep(time.Second)

	base := fmt.Sprintf("/tmp/%d-%d", os.Getpid(), time.Now().UnixNano())
	go runServer(t, "master", "-addr=:6060", "-base="+base, "-vshards=1", "-minGroups=2")

	// NOTE: Here we assume that the master will create the base
	// dir once all required servers register themselves.
	time.Sleep(time.Second)
	_, e := fs.Stat(base)
	assert.Nil(t, e)

}

func runServer(t *testing.T, server string, args ...string) {
	p := path.Join(os.Getenv("GOPATH"), "bin", server)
	fmt.Printf("Starting %v %v\n", p, strings.Join(args, " "))
	b, e := exec.Command(p, args...).CombinedOutput()
	if e != nil {
		t.Skipf("Failed runing %s: %v: %s", server, e, b)
	}
}
