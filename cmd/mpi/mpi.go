package main

import (
	"flag"
	"log"

	mpi "github.com/JohannWeging/go-mpi"
	"github.com/wangkuiyi/phoenix/srvs"
	"github.com/wangkuiyi/phoenix/srvs/aggregator"
	"github.com/wangkuiyi/phoenix/srvs/master"
	"github.com/wangkuiyi/phoenix/srvs/worker"
)

func main() {
	var cfg srvs.Config
	cfg.RegisterFlags()
	masterAddr := flag.String("master", "", "Local master address, must be in form :xxxx")
	timeout := flag.Int("registration", 5, "Registeration timeout in seconds")
	flag.Parse()

	args := flag.Args()
	mpi.Init(&args)
	worldSize, _ := mpi.Comm_size(mpi.COMM_WORLD)
	rank, _ := mpi.Comm_rank(mpi.COMM_WORLD)

	if worldSize < cfg.VShards*(1+cfg.Groups)+1 {
		log.Fatalf("MPI world size %d is less than %d * (1+%d) + 1 = %d", worldSize, cfg.VShards, cfg.Groups, cfg.VShards*(1+cfg.Groups)+1)
	}

	switch {
	case rank == 0:
		master.Run(*masterAddr, *timeout, &cfg)
	case 1 <= rank && rank <= cfg.VShards:
		aggregator.Run(*masterAddr, *timeout)
	default:
		worker.Run(*masterAddr, *timeout)
	}

	mpi.Finalize()
}
