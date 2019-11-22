package traefik

import (
	"io/ioutil"

	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/kit/traefik/internal"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type TraefikConfig struct {
	Endpoint    string `fig:"endpoint"`
	Upstream    string `fig:"upstream"`
	ServiceName string `fig:"service_name"`
	ServicePort string `fig:"service_port"`
}

type Traefik struct {
	disabled bool
	client   internal.Client
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
		client: internal.Client{
			Endpoint: config.Endpoint,
		},
	}
}

func (t *Traefik) WithLog(log *logan.Entry) *Traefik {
	t.log = log
	return t
}

func (t *Traefik) RegisterChi(r chi.Router) error {
	if t.disabled {
		return nil
	}

	configuration, err := t.getConfiguration(r)
	if err != nil {
		return errors.Wrap(err, "failed to get configuration")
	}

	err = t.client.PutConfiguration(configuration)
	if err != nil {
		return errors.Wrap(err, "failed to add configuration")
	}

	return nil
}

func (t *Traefik) getConfiguration(r chi.Router) (internal.Configuration, error) {
	routes, err := internal.NewChi(r).Routes()
	if err != nil {
		return internal.Configuration{}, errors.Wrap(err, "failed to walk chi router")
	}

	routers := make(map[string]*internal.Router, 1)
	routers[t.config.ServiceName] = &internal.Router{
		Service: t.config.ServiceName,
		Rule:    routes,
	}

	services := make(map[string]*internal.Service, 1)
	services[t.config.ServiceName] = &internal.Service{
		LoadBalancer: &internal.ServersLoadBalancer{
			Servers: []internal.Server{
				internal.Server{
					URL:    t.config.Upstream,
					Scheme: "http",
					Port:   t.config.ServicePort,
				},
			},
		},
	}

	configuration := internal.Configuration{
		HTTP: &internal.HTTPConfiguration{
			Routers:  routers,
			Services: services,
		}}

	return configuration, nil
}
