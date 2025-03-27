package stack

import (
	"github.com/akirasalvare/fdu-connect/internal/zcdns"
	"net"
)

type Stack interface {
	Run()
	SetupResolve(r zcdns.LocalServer)
	DialTCP(addr *net.TCPAddr) (net.Conn, error)
	DialUDP(addr *net.UDPAddr) (net.Conn, error)
}
