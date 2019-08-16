package main

import (
	//	"net"
	"flag"
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
	pt.Log(pt.LogSeverityNotice, "Starting samclient")
	s, err = samptc.NewSAMClientPlug()
	if err != nil {
		pt.Log(pt.LogSeverityNotice, err.Error())
		os.Exit(1)
	} else {
		pt.Log(pt.LogSeverityNotice, "samclient ready")
	}
	s.Destination = *Destination
	s.PtInfo, err = pt.ClientSetup(nil)
	if err != nil {
		os.Exit(1)
	}
	s.Run()
	pt.SmethodsDone()
}
