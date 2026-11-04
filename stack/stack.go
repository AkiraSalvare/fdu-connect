package stack

import (
	"context"
	"net"

	"github.com/akirasalvare/fdu-connect/client"
	"github.com/akirasalvare/fdu-connect/internal/ippool"
	"github.com/akirasalvare/fdu-connect/internal/zcdns"
)

type Stack interface {
	Run()
	SetupResolve(r zcdns.LocalServer)
	SetupIPPool(ipPool *ippool.IPPool[[]client.DomainResource])
	DialTCP(ctx context.Context, addr *net.TCPAddr) (net.Conn, error)
	DialUDP(ctx context.Context, addr *net.UDPAddr) (net.Conn, error)
}
