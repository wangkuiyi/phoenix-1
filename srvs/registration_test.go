package srvs

import (
	"os/exec"
	"path"
	"testing"
)

func TestRegistration(t *testing.T) {
	buildServers(t)

	// // NOTE: here we assume only one path in $GOPATH.
	// // BUG: here we assume master port is not being occupied.
	// e = exec.Command(path.Join(os.Getenv("GOPATH"), "worker"), "-master=:6060")
	// if e != nil {
	// 	t.Skipf("Cannot start worker")
	// }

	// e = exec.Command(path.Join(os.Getenv("GOPATH"), "worker"), "-master=:6060")
	// if e != nil {
	// 	t.Skipf("Cannot start worker")
	// }

	// e = exec.Command(path.Join(os.Getenv("GOPATH"), "worker"), "-master=:6060")
	// if e != nil {
	// 	t.Skipf("Cannot start worker")
	// }

}

// server could be "master", "worker", or "aggregator".
func buildServer(server string, t *testing.T) {
	p := path.Join("github.com/wangkuiyi/phoenix/srvs", server)
	e := exec.Command("go", "install", p).Run()
	if e != nil {
		t.Skipf("Failed building %s: %v", server, e)
	}
}

func buildServers(t *testing.T) {
	buildServer("master", t)
	buildServer("worker", t)
	buildServer("aggregator", t)
}
