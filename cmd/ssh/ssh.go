package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/wangkuiyi/phoenix/srvs"
)

func main() {
	var cfg srvs.Config
	cfg.RegisterFlags()
	role := flag.String("role", "", "Process role: master, aggregator or worker")
	addr := flag.String("master", "", "Local master address, must be in form :xxxx")
	timeout := flag.Int("registration", 5, "Registeration timeout in seconds")
	logPrefix := flag.String("logPrefix", "", "Log output file")
	flag.Parse()

	if len(*logPrefix) > 0 {
		logFile := fmt.Sprintf("%s_%s_%d.log", *logPrefix, *role, os.Getpid())
		f, e := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if e != nil {
			log.Panicf("Failed opening log file %s: %v", logFile, e)
		}
		defer f.Close()
		log.SetOutput(f)
		log.SetPrefix(fmt.Sprintf("%s_%d", *role, os.Getpid()))
	}

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
