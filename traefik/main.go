package traefik

import (
	"context"
	"io/ioutil"

	"github.com/go-redis/redis/v8"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type TraefikConfig struct {
	Redis     *redis.Client `fig:"redis"`
	Endpoints []string      `fig:"endpoints"`
}

type Traefik struct {
	disabled bool
	clients  []Client
	log      *logan.Entry
}

func NewNoOp() *Traefik {
	return &Traefik{
		disabled: true,
	}
}

func NewWithRestInit(endpoints []string) *Traefik {
	trfk := Traefik{
		log: logan.New().Out(ioutil.Discard),
	}

	for _, e := range endpoints {
		trfk.clients = append(trfk.clients, &RestClient{
			Endpoint: e,
		})
	}

	return &trfk
}

func NewWithRedisInit(rc *redis.Client) *Traefik {
	return &Traefik{
		log: logan.New().Out(ioutil.Discard),
		clients: []Client{
			&RedisClient{
				ctx: context.TODO(),
				rc:  rc,
			},
		},
	}
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
