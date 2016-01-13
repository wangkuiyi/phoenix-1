package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"testing"
)

func buildBinary() {
	b, e := exec.Command("go", "install", "github.com/wangkuiyi/phoenix/cmd/phoenix-ssh").CombinedOutput()
	if e != nil {
		log.Panicf("Failed building: %v: %s", e, b)
	}
}

func runServer(t *testing.T, role string, args ...string) {
	p := path.Join(GoPath(), "bin", "phoenix-ssh")
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
