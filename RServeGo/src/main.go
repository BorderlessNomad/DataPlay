package Rserve

import (
	"errors"
	"fmt"
	"net"
	"strings"
)

type RServeConnection struct {
	Hello                string
	ServerBanner         string
	AllowUnknownVersions bool
	connection           net.Conn
}

func New() RServeConnection {
	a := RServeConnection{}
	return a
}

func (self RServeConnection) Connect(IP string, port int) error {
	if port > 65536 {
		erm := errors.New("The TCP Stack does not allow ports to be above 2^16")
		return erm
	}
	servAddr := fmt.Sprintf("%s:%d", IP, port)
	tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
	if err != nil {
		erm := errors.New(fmt.Sprintf("Could not resolve %s", servAddr))
		return erm
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		erm := errors.New(fmt.Sprintf("Could not dial %s", servAddr))
		return erm
	}
	buffer := make([]byte, 2048)
	r, e := conn.Read(buffer)
	if e != nil {
		return e
	}
	realbuffer := buffer[:r]
	strbuf := string(realbuffer)
	/*
		Rsrv0103QAP1

		--------------
	*/
	handshakelines := strings.Split(strbuf, "\n")
	self.ServerBanner = handshakelines[0]
	if (len(handshakelines) < 2 && strings.HasPrefix(handshakelines[0], "Rsrv0103QAP1")) || (self.AllowUnknownVersions && strings.HasPrefix(handshakelines[0], "Rsrv")) {
		// Umm, I guess the connection worked then
	} else {
		return fmt.Errorf("Unsupported API version, This could work but I am not going to risk it version: '%s'", handshakelines[0])
	}
	return nil
}

func getcommandcode(method string) byte {

	if method == "eval" {
		return 0x03
	} else if method == "voidEval" {
		return 0x02
	} else if method == "login" {
		return 0x01
	} else {
		// wat.
		fmt.Printf("WARNING, you asked for '%s' thats not a valid command code.", getcommandcode)
		return 8 // Best number to return
	}

}
