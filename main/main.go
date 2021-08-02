package main

import (
	"flag"
	"log"

	"github.com/4gatepylon/GoPoker/net"
)

var runClient *bool = flag.Bool("run-client", true, "Decide whether to run client or server. Default is client (true).")
func main() {
	flag.Parse()
	if runClient == nil {
		log.Fatalf("Must pick client or server\n")
		return
	}
	if *runClient {
		log.Printf("Running client\n")
		net.RunServer()
		return
	}
	net.RunClient()
	log.Printf("Running server\n")
}