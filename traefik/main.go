package traefik

import (
	"io/ioutil"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type TraefikConfig struct {
	Endpoints []string `fig:"endpoints"`
}

type Traefik struct {
	disabled bool
	clients  []Client
	config   TraefikConfig
	log      *logan.Entry
}

func NewNoOp() *Traefik {
	return &Traefik{
		disabled: true,
	}
}

func New(config TraefikConfig) *Traefik {
	trfk := Traefik{
		config: config,
		log:    logan.New().Out(ioutil.Discard),
	}

	for _, e := range config.Endpoints {
		trfk.clients = append(trfk.clients, Client{
			Endpoint: e,
		})
	}

	return &trfk
}

func (t *Traefik) WithLog(log *logan.Entry) *Traefik {
	t.log = log
	return t
}

func (t *Traefik) RegisterConfiguration(configuration Configuration) error {
	if t.disabled {
		return nil
	}

	for _, c := range t.clients {
		err := c.PutConfiguration(configuration)
		if err != nil {
			return errors.Wrap(err, "failed to add configuration")
		}
	}

	return nil
}
