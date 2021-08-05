package main

import (
	"flag"
	"log"
)

func main() {
	var runClient *bool = flag.Bool("client", true, "Decide whether to run client or server. Default is client (true).")

	flag.Parse()
	if runClient == nil {
		log.Fatalf("Must pick client or server\n")
		return
	}
	if *runClient {
		log.Printf("Running client\n")
		return
	}
	log.Printf("Running server\n")
}
