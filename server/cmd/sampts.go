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
	var KeysPath string
	var ClientPath string
	if len(os.Args) > 1 {
		ClientPath = os.Args[1]
	}
	if ClientPath == "" {
		ClientPath = "sam.torrc"
	}
	KeysPath = "sam.torrc.i2pkeys"
	s, err = sampts.NewSAMServerPlug()
	s.PtInfo, err = pt.ServerSetup(nil)
	if err != nil {
		os.Exit(1)
	}
	s.KeysPath = KeysPath
	s.ClientPath = ClientPath
	s.Run()
	pt.SmethodsDone()
}
