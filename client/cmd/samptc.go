package main

import (
	//	"net"
	"flag"
	"log"
	"os"

	"git.torproject.org/pluggable-transports/goptlib.git"
	"github.com/RTradeLtd/sam-pt/client"
)

var s *samptc.SAMClientPlug

var (
	Destination = flag.String("dest", "", "")
)

func main() {
	var err error
	flag.Parse()
	s, err = samptc.NewSAMClientPlug(*Destination)
	if err != nil {
		log.Println("samclient: Client error")
		os.Exit(1)
	}
	s.Destination = *Destination
	s.PtInfo, err = pt.ClientSetup(nil)
	if err != nil {
		os.Exit(1)
	}
	s.Run()
	pt.SmethodsDone()
}
