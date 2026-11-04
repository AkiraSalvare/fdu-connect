package tcptunnel

import (
	"github.com/akirasalvare/fdu-connect/client"
	"github.com/akirasalvare/fdu-connect/internal/ippool"
	"github.com/akirasalvare/fdu-connect/internal/zcdns"
)

type Stack struct {
	client  client.Client
	resolve zcdns.LocalServer
	ipPool  *ippool.IPPool[[]client.DomainResource]
}

func (s *Stack) Run() {}

func NewStack(client client.Client) (*Stack, error) {
	s := &Stack{
		client: client,
	}
	return s, nil
}

func (s *Stack) SetupResolve(r zcdns.LocalServer) {
	s.resolve = r
}

func (s *Stack) SetupIPPool(ipPool *ippool.IPPool[[]client.DomainResource]) {
	s.ipPool = ipPool
}
