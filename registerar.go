package consul

import (
	"github.com/go-kit/kit/sd/consul"
	logger "gitlab.yourypto.com/core/common-modules/logger"

	"github.com/hashicorp/consul/api"
)

// Registrar registers service instance health information to Consul.
type Registrar struct {
	client       consul.Client
	registration *api.AgentServiceRegistration
	logger       logger.Logger
}

// NewRegistrar returns a Consul Registrar acting on the provided catalog
// registration.
func NewRegistrar(client consul.Client, r *api.AgentServiceRegistration) *Registrar {
	return &Registrar{
		client:       client,
		registration: r,
	}
}

// Register implements sd.Registrar interface.
func (p *Registrar) Register() {
	if err := p.client.Register(p.registration); err != nil {
		panic(err)
	}
}

// Deregister implements sd.Registrar interface.
func (p *Registrar) Deregister() {
	_ = p.client.Deregister(p.registration)
}
