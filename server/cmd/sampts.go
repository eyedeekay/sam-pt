package main

import (
	"net"
	"os"

	"git.torproject.org/pluggable-transports/goptlib.git"
	"github.com/RTradeLtd/sam-pt/server"
)

var s *sampts.SAMServerPlug

func Handler(conn net.Conn) error {
	return s.Handler(conn)
}

func AcceptLoop(ln net.Listener) error {
	return s.AcceptLoop(ln)
}

func main() {
	var err error
    if len(os.Args) > 1 {
        s.KeysPath = os.Args[1]
    )else{
        s.KeysPath = "sam.torrc"
    }
	s, err = sampts.NewSAMServerPlug()
	s.PtInfo, err = pt.ServerSetup(nil)
	if err != nil {
		os.Exit(1)
	}
	s.Run()
	pt.SmethodsDone()
}
