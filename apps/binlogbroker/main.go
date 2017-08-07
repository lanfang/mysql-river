package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

var usage = `
Usage of ./binlogbroker:
  -c string
    	json file config
  -etcd string
    	etcd addr
  -g string
    	this parameter is must, group name like mysql instance, ex: db_mall
`

func main() {
	if len(os.Args) == 1 {
		log.Fatal(usage)
		os.Exit(1)
	}
	if err := initSerer(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	s := startServer()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt, syscall.SIGTERM)
	<-sc
	s.Close()
}
