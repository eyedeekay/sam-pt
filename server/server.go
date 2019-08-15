package sampts

import (
	"io"
	"io/ioutil"
	"log"
	"net"
	"sync"

	"git.torproject.org/pluggable-transports/goptlib.git"
	"github.com/eyedeekay/sam3"
	"github.com/eyedeekay/sam3/i2pkeys"
)

type SAMServerPlug struct {
	sam       *sam3.SAM
	KeysPath  string
	Keys      i2pkeys.I2PKeys
	Session   *sam3.StreamSession
	Listener  *sam3.StreamListener
	Client    *sam3.SAMConn
	PtInfo    pt.ServerInfo
	LocalDest string // this must be a full base64 private key
}

func (s *SAMServerPlug) TorRCClient() string {
	return `UseBridges 1
Bridge sam ` + s.Keys.Addr().Base64() + `

ClientTransportPlugin sam exec /usr/bin/samclient ` + s.Keys.Addr().Base64() + `
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
				s.KeysPath,
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
