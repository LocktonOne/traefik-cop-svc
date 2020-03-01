package traefik

import (
	"io/ioutil"
	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type TraefikConfig struct {
	Endpoint                   string        `fig:"endpoint"`
	GeneralHealthcheckTimeout  time.Duration `fig:"general_healthcheck_timeout"`
	GeneralHealthcheckInterval time.Duration `fig:"general_healthcheck_interval"`
}

type Traefik struct {
	disabled bool
	client   Client
	config   TraefikConfig
	log      *logan.Entry
}

func NewNoOp() *Traefik {
	return &Traefik{
		disabled: true,
	}
}

func New(config TraefikConfig) *Traefik {
	return &Traefik{
		config: config,
		log:    logan.New().Out(ioutil.Discard),
		client: Client{
			Endpoint: config.Endpoint,
		},
	}
}

func (c *Traefik) Cfg() *TraefikConfig {
	return &c.config
}

func (t *Traefik) WithLog(log *logan.Entry) *Traefik {
	t.log = log
	return t
}

func (t *Traefik) RegisterConfiguration(configuration Configuration) error {
	if t.disabled {
		return nil
	}

	err := t.client.PutConfiguration(configuration)
	if err != nil {
		return errors.Wrap(err, "failed to add configuration")
	}

	return nil
}
