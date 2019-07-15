package client

import (
	"fmt"

	"github.com/hashicorp/consul/api"
)

// Client is consul client wrapper
type Client interface {
	// CreateService creates consul service
	CreateService(Service, bool, *api.QueryOptions) ([]*api.ServiceEntry, *api.QueryMeta, error)
	// Register a service with local agent
	Register(Service) error
	// DeRegister a service with consul agent
	DeRegister(string) error
}

type client struct {
	consul *api.Client
}

// Register a service with consul agent
func (c *client) Register(service Service) error {
	registry := &api.AgentServiceRegistration{
		ID:   service.Name,
		Name: service.Name,
		Port: service.Port,
	}
	return c.consul.Agent().ServiceRegister(registry)
}

// NewConsulClientFromConfig creates consul client using consul client
func NewConsulClientFromConfig(config *api.Config) (Client, error) {
	c, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &client{consul: c}, nil
}

// NewConsulClient creates consul client if addr is passed uses that address else uses default
func NewConsulClient(addr string) (Client, error) {
	config := api.DefaultConfig()
	if addr != "" {
		config.Address = addr
	}
	return NewConsulClientFromConfig(config)
}

// DeRegister a service with consul agent
func (c *client) DeRegister(id string) error {
	return c.consul.Agent().ServiceDeregister(id)
}

// CreateService creates consul service
func (c *client) CreateService(service Service, passingOnly bool, queryOptions *api.QueryOptions) ([]*api.ServiceEntry, *api.QueryMeta, error) {
	addrs, meta, err := c.consul.Health().Service(service.Name, service.Tag, passingOnly, queryOptions)
	if len(addrs) == 0 && err == nil {
		return nil, nil, fmt.Errorf("service ( %s ) was not found", service.Name)
	}
	if err != nil {
		return nil, nil, err
	}
	return addrs, meta, nil
}
