package samptc

import (
	"log"
	"strconv"
	"strings"
)

import (
	"git.torproject.org/pluggable-transports/goptlib.git"
	"github.com/eyedeekay/sam-forwarder/hashhash"
	"github.com/eyedeekay/sam3"
	"github.com/eyedeekay/sam3/i2pkeys"
)

var Options_Short = []string{"inbound.length=1", "outbound.length=1",
	"inbound.lengthVariance=0", "outbound.lengthVariance=0",
	"inbound.backupQuantity=1", "outbound.backupQuantity=1",
	"inbound.quantity=2", "outbound.quantity=2", "i2cp.dontPublishLeaseSet=true"}

func NewSAMClientPlug(Destination, BridgeURL string) (*SAMClientPlug, error) {
	var s SAMClientPlug
	var err error
	s.Destination = Destination
	s.BridgeURL = BridgeURL
	s.sam, err = sam3.NewSAM("127.0.0.1:7656")
	if err != nil {
		return nil, err
	}
	log.Println("samclient: SAM connected")
	s.keys, err = s.sam.NewKeys()
	if err != nil {
		return nil, err
	}
	log.Println("samclient: Keys generated")
	s.Session, err = s.sam.NewStreamSessionWithSignature("sam-pt", s.keys, Options_Short, sam3.Sig_EdDSA_SHA512_Ed25519)
	if err != nil {
		return nil, err
	}
	log.Println("samclient: Session started")
	if strings.HasSuffix(s.Destination, ".i2p") {
		log.Println("samclient: I2P Domain Detected", s.Destination)
		s.destaddr, err = s.sam.Lookup(s.Destination)
		if err != nil {
			return nil, err
		}
		log.Println("samclient: Looked up destination")
	} else if slice := strings.Split(s.Destination, " "); len(slice) > 30 {
		if length, err := strconv.Atoi(slice[len(slice)-1]); err == nil {
			Hasher, err := hashhash.NewHasher(length)
			if err != nil {
				return nil, err
			}
			log.Println("samclient: Created hash decoder")
			s.Destination, err = Hasher.Unfriendlyslice(slice[0 : len(slice)-2])
			if err != nil {
				return nil, err
			}
			log.Println("samclient: decoded hash")
		} else {
			Hasher, err := hashhash.NewHasher(52)
			if err != nil {
				return nil, err
			}
			log.Println("samclient: created hash decoder")
			s.Destination, err = Hasher.Unfriendlyslice(slice)
			if err != nil {
				return nil, err
			}
			log.Println("samclient: decoded hash")
		}
		s.destaddr, err = s.sam.Lookup(s.Destination)
		if err != nil {
			return nil, err
		}
		log.Println("samclient: Looked up destination")
	} else {
		s.destaddr, err = i2pkeys.NewI2PAddrFromString(s.Destination)
		if err != nil {
			return nil, err
		}
		log.Println("samclient: created address")
	}
	log.Println("samclient: dialing server")
	s.Client, err = s.Session.DialI2P(s.destaddr)
	if err != nil {
		return nil, err
	}
	log.Println("samclient: dialed address")
	s.PtInfo, err = pt.ClientSetup(nil)
	if err != nil {
		return nil, err
	}
	log.Println("samclient: set up client")
	return &s, nil
}
