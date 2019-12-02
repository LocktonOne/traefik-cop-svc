package cop

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type CopConfig struct {
	Endpoint    string `fig:"endpoint"`
	Upstream    string `fig:"upstream"`
	ServiceName string `fig:"service_name"`
	ServicePort string `fig:"service_port"`
}

type Cop struct {
	disabled bool
	client   Client
	config   CopConfig
	log      *logan.Entry
}

func NewNoOp() *Cop {
	return &Cop{
		disabled: true,
	}
}

func New(config CopConfig) *Cop {
	return &Cop{
		config: config,
		log:    logan.New().Out(ioutil.Discard),
		client: Client{
			Endpoint: config.Endpoint,
		},
	}
}

func (c *Cop) WithLog(log *logan.Entry) *Cop {
	c.log = log
	return c
}

func (c *Cop) RegisterChi(r chi.Router) error {
	if c.disabled {
		return nil
	}

	rule, err := GetRule(r)
	if err != nil {
		return errors.Wrap(err, "failed to get rule")
	}

	service := Service{Data: ServiceData{
		ID:   c.config.ServiceName,
		Type: "traefik-service",
		Attributes: ServiceAttributes{
			Name: c.config.ServiceName,
			Url:  c.config.Upstream,
			Port: c.config.ServicePort,
			Rule: rule,
		},
	}}

	err = c.client.AddService(service)
	if err != nil {
		return errors.Wrap(err, "failed to add service")
	}

	return nil
}

func GetRule(r chi.Router) (string, error) {
	var routes []string
	walk := func(method, route string, _ http.Handler, _ ...func(http.Handler) http.Handler) error {
		route = strings.Replace(route, "/*/", "/", -1)
		if len(route) > 1 {
			route = strings.TrimRight(route, "/")
		}

		routes = append(routes, fmt.Sprintf("Path(`%s`)", route))
		return nil
	}

	err := chi.Walk(r, walk)
	if err != nil {
		return "", errors.Wrap(err, "failed to walk router")
	}
	return strings.Join(routes, "||"), nil
}
