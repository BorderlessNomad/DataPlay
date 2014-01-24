package Rserve

import (
	"errors"
	"fmt"
	"net"
	"strings"
)

type RServeConnection struct {
	Hello      string
	connection net.Conn
}

func New() RServeConnection {
	a := RServeConnection{}
	return a
}

func (RServeConnection) Connect(IP string, port int) error {
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
	if handshakelines[0] == "Rsrv0103QAP1" {
		fmt.Println(":D")
	} else {
		return fmt.Errorf("Unsupported API version, This could work but I am not going to risk it")
	}
	return nil
}
