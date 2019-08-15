package sampts

import (
	"os"
)

import (
	"github.com/eyedeekay/sam3"
	"github.com/eyedeekay/sam3/i2pkeys"
)

var Options_Short = []string{"inbound.length=1", "outbound.length=1",
	"inbound.lengthVariance=0", "outbound.lengthVariance=0",
	"inbound.backupQuantity=1", "outbound.backupQuantity=1",
	"inbound.quantity=2", "outbound.quantity=2"}

func NewSAMServerPlug() (*SAMServerPlug, error) {
	var s SAMServerPlug
	var err error
	s.sam, err = sam3.NewSAM("127.0.0.1:7656")
	if err != nil {
		return nil, err
	}
	if s.KeysPath == "" {
		s.KeysPath = "sam.torrc.i2pkeys"
	}
	if _, err := os.Stat(s.KeysPath); os.IsNotExist(err) {
		s.Keys, err = s.sam.NewKeys()
		if err != nil {
			return nil, err
		}
		file, err := os.Create(s.KeysPath)
		if err != nil {
			return nil, err
		}
		err = i2pkeys.StoreKeysIncompat(s.Keys, file)
	} else {
		file, err := os.Open(s.KeysPath)
		if err != nil {
			return nil, err
		}
		s.Keys, err = i2pkeys.LoadKeysIncompat(file)
		if err != nil {
			return nil, err
		}
	}
	s.Session, err = s.sam.NewStreamSession("sam-pt", s.Keys, Options_Short)
	if err != nil {
		return nil, err
	}
	return &s, nil
}
