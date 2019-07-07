package main

import (
	"github.com/RTradeLtd/sam-pt/server"
)

var pt sampts.SAMServerPlug

func main() {
	if pt, err := sampts.NewSAMServerPlug(); err != nil {
		panic(err)
	} else {
		if err := pt.Run(); err != nil {
			panic(err)
		}
	}
}
