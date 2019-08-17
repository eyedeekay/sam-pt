package main

import (
	"flag"
	"log"
	"os"

	"github.com/RTradeLtd/sam-pt/client"
)

var s *samptc.SAMClientPlug

var (
	Destination = flag.String("destination", "fZAVcqVFgdsalFlDG52ItFwOblIyddiENrrKJ4Fx46qC2AyKJbId9P2l3g-RCUzaWVUge~WEyxb5IqfHt9P8HKvG6cXnUK2sCjlEl-hgbVnyNZVoQzoFD7g8CGtcjIdjXtJby2QYenOdv-Q3uuu44MEVGe4rnzteg85nhUvJ5jPIDZHIj4s5kp4hs0l9KQ~SCpPG4fCBKQCsup26tBOpdk8EwIKJ8nNeeLHTh~ACEbbja0BMWnow8BY3siB936-TTCo~F37SMPP4-H18UPnzAQHThX1yeb9kjsD6EExVJqmeuXmh0ciRDdTredm3wC4ftKDvUqVa4jA8C8WXKqwNKYFErcR2eAhAyLPk66uetVL6IJFci9KM1XzxyO6Dlb7RouTh8WFkKg0TT6NmFjRhQNh9NVruv2oJCpoBNG2krp0tvurAKXQUC-BtA8JR-V880IObwRgYMSStfPTxTrnCBazc6~kQNYJWxQaATHWYQCfExIKkne~02k05kVTWPB0WAAAA", "An i2p address to use as a Tor Bridge")
	Bridge      = flag.String("bridge-tunnel", "127.0.0.1:7951", "TCP address to host the bridge on")
)

func main() {
	var err error
	flag.Parse()
	s, err = samptc.NewSAMClientPlug(*Destination, *Bridge)
	if err != nil {
		log.Println("samclient: Client error", err)
		os.Exit(1)
	}
	log.Println("samclient: Client starting")
	if err = s.Run(); err != nil {
		os.Exit(1)
	}
}
