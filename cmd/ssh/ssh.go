package main

import (
	"flag"
	"log"

	"github.com/wangkuiyi/phoenix/srvs"
	"github.com/wangkuiyi/phoenix/srvs/aggregator"
	"github.com/wangkuiyi/phoenix/srvs/master"
	"github.com/wangkuiyi/phoenix/srvs/worker"
)

func main() {
	var cfg srvs.Config
	cfg.RegisterFlags()
	role := flag.String("role", "", "Process role: master, aggregator or worker")
	addr := flag.String("master", "", "Local master address, must be in form :xxxx")
	timeout := flag.Int("registration", 5, "Registeration timeout in seconds")
	flag.Parse()

	switch *role {
	case "master":
		master.Run(*addr, *timeout, &cfg)
	case "aggregator":
		aggregator.Run(*addr, *timeout)
	case "worker":
		worker.Run(*addr, *timeout)
	default:
		log.Fatal("Unknown role: ", *role)
	}
}
