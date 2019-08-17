package samptc

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"

	"git.torproject.org/pluggable-transports/goptlib.git"
	"github.com/eyedeekay/sam3"
	"github.com/eyedeekay/sam3/i2pkeys"
)

type SAMClientPlug struct {
	sam         *sam3.SAM
	keys        i2pkeys.I2PKeys
	destaddr    i2pkeys.I2PAddr
	Session     *sam3.StreamSession
	Client      *sam3.SAMConn
	PtInfo      pt.ClientInfo
	Destination string
	BridgeURL   string
}

func (s *SAMClientPlug) NetworkListener() net.Listener {
	listener, _ := s.Session.Listen()
	return listener
}

func (s *SAMClientPlug) TransportAcceptI2P() (*SAMClientPlug, error) {
	return s, nil
}

func (s *SAMClientPlug) Close() error {
	return s.Close()
}

func (s *SAMClientPlug) CopyLoop(or net.Conn) {
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

func (s *SAMClientPlug) Handler(conn *pt.SocksConn) error {
	defer conn.Close()
	err := conn.Grant(nil)
	if err != nil {
		return err
	}
	// do something with conn and remote.
	s.CopyLoop(conn)
	return nil
}

func (s *SAMClientPlug) AcceptLoop(ln *pt.SocksListener) error {
	defer ln.Close()
	for {
		conn, err := ln.AcceptSocks()
		if err != nil {
			if e, ok := err.(net.Error); ok && e.Temporary() {
				pt.Log(pt.LogSeverityError, "accept error: "+err.Error())
				continue
			}
			return err
		}
		go s.Handler(conn)
	}
	return nil
}

func (s *SAMClientPlug) Run() error {
	var err error
	log.Println("samclient: dialing server")
	s.Client, err = s.Session.DialI2P(s.destaddr)
	if err != nil {
		return err
	}
	log.Println("samclient: dialed address")
	if s.PtInfo.ProxyURL != nil {
		// you need to interpret the proxy URL yourself
		// call pt.ProxyDone instead if it's a type you understand
		pt.ProxyError(fmt.Sprintf("proxy %s is not supported", s.PtInfo.ProxyURL))
		return fmt.Errorf("proxy %s is not supported", s.PtInfo.ProxyURL)
	}
	for _, methodName := range s.PtInfo.MethodNames {
		switch methodName {
		case "sam":
			ln, err := pt.ListenSocks("tcp", "127.0.0.1:0")
			if err != nil {
				pt.CmethodError(methodName, err.Error())
				break
			}
			go s.AcceptLoop(ln)
			pt.Cmethod(methodName, ln.Version(), ln.Addr())
		default:
			pt.CmethodError(methodName, "no such method")
		}
	}
	pt.CmethodsDone()
	return nil
}
