package main

import (
	"flag"
	"log"

	"github.com/wangkuiyi/phoenix/srvs"
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
		srvs.RunMaster(*addr, *timeout, &cfg)
	case "aggregator":
		srvs.RunAggregator(*addr, *timeout)
	case "worker":
		srvs.RunWorker(*addr, *timeout)
	default:
		log.Fatal("Unknown role: ", *role)
	}
}
