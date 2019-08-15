package sampts

import (
	"io"
	"io/ioutil"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"

	"git.torproject.org/pluggable-transports/goptlib.git"
	"github.com/eyedeekay/sam-forwarder/hashhash"
	"github.com/eyedeekay/sam3"
	"github.com/eyedeekay/sam3/i2pkeys"
)

type SAMServerPlug struct {
	sam        *sam3.SAM
	KeysPath   string
	ClientPath string
	Keys       i2pkeys.I2PKeys
	Session    *sam3.StreamSession
	Listener   *sam3.StreamListener
	Client     *sam3.SAMConn
	PtInfo     pt.ServerInfo
}

func (s *SAMServerPlug) TorRCClient() string {
	return `
## Conflgure a client by adding these lines to your torrc

UseBridges 1
Bridge sam ` + s.Keys.Addr().Base64() + `

ClientTransportPlugin sam exec /usr/bin/samclient ` + s.Keys.Addr().Base64() + `

## OR you can use the base32

UseBridges 1
Bridge sam ` + s.Keys.Addr().Base32() + `

ClientTransportPlugin sam exec /usr/bin/samclient ` + s.Keys.Addr().Base32() + `

## OR you can use a readable mnemonic

UseBridges 1
Bridge sam "` + s.Mnemonic() + " " + strconv.Itoa(len(strings.Replace(s.Keys.Addr().Base32(), ".b32.i2p", "", 1))) + `"

ClientTransportPlugin sam exec /usr/bin/samclient "` + s.Mnemonic() + " " + strconv.Itoa(len(strings.Replace(s.Keys.Addr().Base32(), ".b32.i2p", "", 1))) + `"
`
}

func (s *SAMServerPlug) NetworkListener() net.Listener {
	listener, _ := s.Session.Listen()
	return listener
}

func (s *SAMServerPlug) TransportAcceptI2P() (*SAMServerPlug, error) {
	return s, nil
}

func (s *SAMServerPlug) Close() error {
	return s.Close()
}

func (s *SAMServerPlug) CopyLoop(or net.Conn) {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		io.Copy(or, s.Client)
		wg.Done()
	}()
	go func() {
		io.Copy(s.Client, or)
		wg.Done()
	}()

	wg.Wait()
}

func (s *SAMServerPlug) Handler(conn net.Conn) error {
	defer conn.Close()
	or, err := pt.DialOr(&s.PtInfo, conn.RemoteAddr().String(), "foo")
	if err != nil {
		return err
	}
	defer or.Close()
	// do something with or and conn
	s.CopyLoop(or)
	return nil
}

func (s *SAMServerPlug) AcceptLoop(ln net.Listener) error {
	defer ln.Close()
	var err error
	s.Client, err = ln.(*sam3.StreamListener).AcceptI2P()
	for {
		if err != nil {
			if e, ok := err.(net.Error); ok && e.Temporary() {
				continue
			}
			pt.Log(pt.LogSeverityError, "accept error: "+err.Error())
			return err
		}
		go s.Handler(s.Client)
	}
	return nil
}

func (s *SAMServerPlug) Mnemonic() string {
	b32 := strings.Replace(s.Keys.Addr().Base32(), ".b32.i2p", "", 1)
	Hasher, err := hashhash.NewHasher(len(b32))
	if err != nil {
		return s.Keys.Addr().Base32()
	}
	hash, _ := Hasher.Friendly(b32)
	return hash
}

func (s *SAMServerPlug) Run() error {
	var err error
	for _, bindaddr := range s.PtInfo.Bindaddrs {
		switch bindaddr.MethodName {
		case "samserver":
			s.Listener, err = s.Session.Listen()
			if err != nil {
				pt.SmethodError(bindaddr.MethodName, err.Error())
				break
			}
			log.Println(s.TorRCClient())
			ioutil.WriteFile(
				s.ClientPath,
				[]byte(s.TorRCClient()),
				0644,
			)
			go s.AcceptLoop(s.Listener)
			pt.Smethod(bindaddr.MethodName, s.Listener.Addr())
		default:
			pt.SmethodError(bindaddr.MethodName, "no such method")
		}
	}
	pt.SmethodsDone()
	return nil
}
