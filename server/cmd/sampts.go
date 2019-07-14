package main

import (
	"io"
	"net"
	"os"
	"sync"

	"git.torproject.org/pluggable-transports/goptlib.git"
	"github.com/RTradeLtd/sam-pt/server"
	"github.com/eyedeekay/goSam"
)

var s sampts.SAMServerPlug

func copyLoop(stream *goSam.Client, or net.Conn) {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		io.Copy(or, stream)
		wg.Done()
	}()
	go func() {
		io.Copy(stream, or)
		wg.Done()
	}()

	wg.Wait()
}

func Handler(conn net.Conn) error {
	defer conn.Close()
	or, err := pt.DialOr(&s.PtInfo, conn.RemoteAddr().String(), "foo")
	if err != nil {
		return err
	}
	defer or.Close()
	// do something with or and conn
	copyLoop(s.Client, or)
	return nil
}

func AcceptLoop(ln net.Listener) error {
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
		go Handler(conn)
	}
	return nil
}

func main() {
	var err error
	s.PtInfo, err = pt.ServerSetup(nil)
	if err != nil {
		os.Exit(1)
	}
	for _, bindaddr := range s.PtInfo.Bindaddrs {
		switch bindaddr.MethodName {
		case "samserver":
			ln, err := s.Client.ListenI2P(s.LocalDest)
			if err != nil {
				pt.SmethodError(bindaddr.MethodName, err.Error())
				break
			}
			go AcceptLoop(ln)
			pt.Smethod(bindaddr.MethodName, ln.Addr())
		default:
			pt.SmethodError(bindaddr.MethodName, "no such method")
		}
	}
	pt.SmethodsDone()
}
