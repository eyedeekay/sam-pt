package main

import (
	"github.com/RTradeLtd/sam-pt/client"
)

func main() {
	if pt, err := samptc.NewSAMClientPlug(); err != nil {
		panic(err)
	} else {
		if err := pt.Run(); err != nil {
			panic(err)
		}
	}
}
