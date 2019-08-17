package sampts

import (
	"log"
	"os"
)

import (
	"github.com/eyedeekay/sam3"
	"github.com/eyedeekay/sam3/i2pkeys"
)

var Options_Short = []string{"inbound.length=1", "outbound.length=1",
	"inbound.lengthVariance=0", "outbound.lengthVariance=0",
	"inbound.backupQuantity=1", "outbound.backupQuantity=1",
	"inbound.quantity=2", "outbound.quantity=2", "i2cp.dontPublishLeaseSet=false"}

func NewSAMServerPlug(KeysPath, ClientPath string) (*SAMServerPlug, error) {
	var s SAMServerPlug
	var err error
	s.sam, err = sam3.NewSAM("127.0.0.1:7656")
	if err != nil {
		return nil, err
	}
	s.KeysPath = KeysPath
	s.ClientPath = ClientPath
	log.Println("samserver: Connected to SAM")
	log.Println("samserver: Checking for keys in" + s.KeysPath)
	if _, err := os.Stat(s.KeysPath); os.IsNotExist(err) {
		s.Keys, err = s.sam.NewKeys()
		if err != nil {
			return nil, err
		}
		log.Println("samserver: Generated keys")
		file, err := os.Create(s.KeysPath)
		if err != nil {
			return nil, err
		}
		err = i2pkeys.StoreKeysIncompat(s.Keys, file)
		if err != nil {
			return nil, err
		}
		log.Println("samserver: Saved keys")
	} else {
		file, err := os.Open(s.KeysPath)
		if err != nil {
			return nil, err
		}
		log.Println("samserver: Found keys")
		s.Keys, err = i2pkeys.LoadKeysIncompat(file)
		if err != nil {
			return nil, err
		}
		log.Println("samserver: Loaded Keys")
	}
	return &s, nil
}
