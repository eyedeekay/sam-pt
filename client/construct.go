package samptc

import (
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
	s.destaddr, err = i2pkeys.NewI2PAddrFromString(s.destination)
	if err != nil {
		return nil, err
	}
	s.Client, err = s.Session.DialI2P(s.destaddr)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// NewSAMClientPlugFromOptions creates a new client, connecting to a specified port
/*func NewSAMClientPlugFromOptions(opts ...func(*goSam.Client) error) (*SAMClientPlug, error) {
	var c SAMClientPlug
	for _, o := range opts {
		if err := o(c.Client); err != nil {
			return nil, err
		}
	}
	return &c, nil
}*/
