package sampts

import (
	"github.com/eyedeekay/sam3"
	//	"github.com/eyedeekay/sam3/i2pkeys"
)

var Options_Short = []string{"inbound.length=1", "outbound.length=1",
	"inbound.lengthVariance=0", "outbound.lengthVariance=0",
	"inbound.backupQuantity=1", "outbound.backupQuantity=1",
	"inbound.quantity=2", "outbound.quantity=2"}

func NewSAMServerPlug() (*SAMServerPlug, error) {
	var s SAMServerPlug
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
	return &s, nil
}

// NewSAMClientPlugFromOptions creates a new client, connecting to a specified port
/*func NewSAMServerPlugFromOptions(opts ...func(*goSam.Client) error) (*SAMServerPlug, error) {
	var c SAMServerPlug
	for _, o := range opts {
		if err := o(c.Client); err != nil {
			return nil, err
		}
	}
	return &c, nil
}*/
