package main

import (
	//	"net"
	"os"

	"git.torproject.org/pluggable-transports/goptlib.git"
	"github.com/RTradeLtd/sam-pt/client"
)

var s *samptc.SAMClientPlug

func Handler(conn *pt.SocksConn) error {
	return s.Handler(conn)
}

func AcceptLoop(ln *pt.SocksListener) error {
	return s.AcceptLoop(ln)
}

func main() {
	var err error
	s.Destination = os.Args[1]
	s, err = samptc.NewSAMClientPlug()
	s.PtInfo, err = pt.ClientSetup(nil)
	if err != nil {
		os.Exit(1)
	}
	s.Run()
	pt.SmethodsDone()
}
