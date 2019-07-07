package samptc

import (
	"fmt"
	"net"
	"os"

	"git.torproject.org/pluggable-transports/goptlib.git"
	"github.com/eyedeekay/goSam"
)

type SAMClientPlug struct {
	*goSam.Client
	ptInfo      pt.ClientInfo
	destination string
}

func (s *SAMClientPlug) NetworkListener() net.Listener {
	listener, _ := s.Listen()
	return listener
}

func (s *SAMClientPlug) TransportAcceptI2P() (*SAMClientPlug, error) {
	return s, nil
}

func (s *SAMClientPlug) Close() error {
	return s.Close()
}

func (s *SAMClientPlug) Handler(conn *pt.SocksConn) error {
	defer conn.Close()
	remote, err := s.Client.Dial("i2p", s.destination) //conn.Req.Target)
	if err != nil {
		conn.Reject()
		return err
	}
	defer remote.Close()
	err = conn.Grant(remote.RemoteAddr().(*net.TCPAddr))
	if err != nil {
		return err
	}
	// do something with conn and remote.
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

func (s *SAMClientPlug) Run() {
	var err error
	s.ptInfo, err = pt.ClientSetup(nil)
	if err != nil {
		os.Exit(1)
	}
	if s.ptInfo.ProxyURL != nil {
		// you need to interpret the proxy URL yourself
		// call pt.ProxyDone instead if it's a type you understand
		pt.ProxyError(fmt.Sprintf("proxy %s is not supported", s.ptInfo.ProxyURL))
		os.Exit(1)
	}
	for _, methodName := range s.ptInfo.MethodNames {
		switch methodName {
		case "samclient":
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
}
