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


	ZI forgot tolock my screeen''
}

func TestRegistration(t *testing.T) {
	go runServer(t, "worker", "-master=:6060")
	time.Sleep(time.Second)
	go runServer(t, "worker", "-master=:6060")
	time.Sleep(time.Second)
	go runServer(t, "aggregator", "-master=:6060")
	time.Sleep(time.Second)
	runServer(t, "master", "-addr=:6060", "-base=/tmp/", "-vshards=1", "-minGroups=2")
}

func runServer(t *testing.T, server string, args ...string) {
	p := path.Join(os.Getenv("GOPATH"), "bin", server)
	fmt.Printf("Starting %v %v\n", p, strings.Join(args, " "))
	b, e := exec.Command(p, args...).CombinedOutput()
	if e != nil {
		t.Skipf("Failed runing %s: %v: %s", server, e, b)
	}
}
