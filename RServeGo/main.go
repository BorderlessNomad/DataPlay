package Rserve

import (
	"fmt"
	_ "net"
)

type RServeConnection struct {
	Hello string
}

func (RServeConnection) Connect() {
	fmt.Println("Hai")
}

func New() RServeConnection {
	a := RServeConnection{}
	return a
}
