package sampts

import (
	"git.torproject.org/pluggable-transports/goptlib.git"
	"github.com/eyedeekay/goSam"
	"net"
)

type SAMServerPlug struct {
	*goSam.Client
	destination string
	ptInfo      pt.ServerInfo
}

func (s *SAMServerPlug) NetworkListener() net.Listener {
	listener, _ := s.Listen()
	return listener
}

func (s *SAMServerPlug) TransportAcceptI2P() (*SAMServerPlug, error) {
	return s, nil
}

func (s *SAMServerPlug) Close() error {
	return s.Close()
}

func (s *SAMServerPlug) Handler(conn net.Conn) error {
	defer conn.Close()
	or, err := pt.DialOr(&s.ptInfo, conn.RemoteAddr().String(), "foo")
	if err != nil {
		return err
	}
	defer or.Close()
	// do something with or and conn
	return nil
}

func (s *SAMServerPlug) AcceptLoop(ln net.Listener) error {
	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			if e, ok := err.(net.Error); ok && e.Temporary() {
				continue
			}
			pt.Log(pt.LogSeverityError, "accept error: "+err.Error())
			return err
		}
		go s.Handler(conn)
	}
	return nil
}
