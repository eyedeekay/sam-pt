package sampts

import (
	"io"
	"io/ioutil"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
)
import (
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
Bridge sam 127.0.0.1:7951

ClientTransportPlugin sam exec /usr/bin/samclient -destination=` + s.Keys.Addr().Base64() + `

## OR you can use the base32
UseBridges 1
Bridge sam 127.0.0.1:7951

ClientTransportPlugin sam exec /usr/bin/samclient -destination=` + s.Keys.Addr().Base32() + `

## OR you can use a readable mnemonic
UseBridges 1
Bridge sam 127.0.0.1:7951

ClientTransportPlugin sam exec /usr/bin/samclient -destination="` + s.Mnemonic() + " " + strconv.Itoa(len(strings.Replace(s.Keys.Addr().Base32(), ".b32.i2p", "", 1))) + `"
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
	or, err := pt.DialOr(&s.PtInfo, conn.RemoteAddr().String(), "sam")
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
	for {
		acc, err := ln.Accept()
		if err != nil {
			return err
		}
		s.Client, err = s.Listener.AcceptI2P()
		if err != nil {
			return err
		}
		if err != nil {
			if e, ok := err.(net.Error); ok && e.Temporary() {
				continue
			}
			pt.Log(pt.LogSeverityError, "accept error: "+err.Error())
			return err
		}
		go s.Handler(acc)
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
	log.Println("## CONFIGURE samserver CLIENT " + s.TorRCClient())
	if err := ioutil.WriteFile(
		s.ClientPath,
		[]byte(s.TorRCClient()),
		0644,
	); err != nil {
		log.Println("samserver: Couldn't write config to file" + s.ClientPath)
		return err
	}
	var err error
	s.Session, err = s.sam.NewStreamSession("sam-pt", s.Keys, Options_Short)
	if err != nil {
		return err
	}
	log.Println("samserver: SAM Session created")
	s.Listener, err = s.Session.Listen()
	if err != nil {
		return err
	}
	log.Println("samserver: SAM Listener created")
	s.PtInfo, err = pt.ServerSetup(nil)
	if err != nil {
		return err
	}
	log.Println("samserver: SAM Pluggable Transport Configured")
	for _, bindaddr := range s.PtInfo.Bindaddrs {
		switch bindaddr.MethodName {
		case "sam":
			ln, err := net.ListenTCP("tcp", bindaddr.Addr)
			if err != nil {
				pt.SmethodError(bindaddr.MethodName, err.Error())
				break
			}
			go s.AcceptLoop(ln)
			pt.Smethod(bindaddr.MethodName, ln.Addr())
		default:
			pt.SmethodError(bindaddr.MethodName, "no such method")
		}
	}
	pt.SmethodsDone()
	return nil
}
