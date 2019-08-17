package main

import (
	"flag"
	"os"

	"git.torproject.org/pluggable-transports/goptlib.git"
	"github.com/RTradeLtd/sam-pt/server"
)

var s *sampts.SAMServerPlug

var (
	ClientPath = flag.String("client-config", "sam.torrc", "Create client config examples here")
	KeysPath   = flag.String("i2p-keys", "sam.torrc.i2pkeys", "Create I2P keys here for storage")
)

func main() {
	var err error
	flag.Parse()
	s, err = sampts.NewSAMServerPlug(*KeysPath, *ClientPath)
	if err != nil {
		pt.Log(pt.LogSeverityError, err.Error())
		os.Exit(1)
	}
	if err := s.Run(); err != nil {
		panic(err)
	}

}
