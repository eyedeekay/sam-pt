package samptc

import (
	"git.torproject.org/pluggable-transports/goptlib.git"
	"github.com/eyedeekay/goSam"
	"net"
)

type SAMClientPlug struct {
	*goSam.Client
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

func NewSAMClientPlug() (*SAMClientPlug, error) {
	var s SAMClientPlug
	var err error
	s.Client, err = goSam.NewDefaultClient()
	if err != nil {
		return nil, err
	}
	return &s, nil
}
