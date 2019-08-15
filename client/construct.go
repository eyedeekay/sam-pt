package samptc

import (
	"strconv"
	"strings"
)

import (
	"github.com/eyedeekay/sam-forwarder/hashhash"
	"github.com/eyedeekay/sam3"
	"github.com/eyedeekay/sam3/i2pkeys"
)

var Options_Short = []string{"inbound.length=1", "outbound.length=1",
	"inbound.lengthVariance=0", "outbound.lengthVariance=0",
	"inbound.backupQuantity=1", "outbound.backupQuantity=1",
	"inbound.quantity=2", "outbound.quantity=2"}

func NewSAMClientPlug() (*SAMClientPlug, error) {
	var s SAMClientPlug
	var err error
	s.keys, err = s.sam.NewKeys()
	if err != nil {
		return nil, err
	}
	s.sam, err = sam3.NewSAM("127.0.0.1:7656")
	if err != nil {
		return nil, err
	}
	s.Session, err = s.sam.NewStreamSession("sam-pt", s.keys, Options_Short)
	if err != nil {
		return nil, err
	}
	if strings.HasSuffix(s.Destination, ".i2p") {
		s.destaddr, err = s.sam.Lookup(s.Destination)
	} else if slice := strings.Split(s.Destination, " "); len(slice) > 30 {
		if length, err := strconv.Atoi(slice[len(slice)-1]); err == nil {
			Hasher, err := hashhash.NewHasher(length)
			if err != nil {
				return nil, err
			}
			s.Destination, err = Hasher.Unfriendlyslice(slice[0 : len(slice)-2])
			if err != nil {
				return nil, err
			}
		} else {
			Hasher, err := hashhash.NewHasher(52)
			if err != nil {
				return nil, err
			}
			s.Destination, err = Hasher.Unfriendlyslice(slice)
			if err != nil {
				return nil, err
			}
		}
		s.destaddr, err = s.sam.Lookup(s.Destination)
	} else {
		s.destaddr, err = i2pkeys.NewI2PAddrFromString(s.Destination)
	}
	if err != nil {
		return nil, err
	}
	s.Client, err = s.Session.DialI2P(s.destaddr)
	if err != nil {
		return nil, err
	}
	return &s, nil
}
