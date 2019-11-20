package janus

import (
	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/kit/janus/internal"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

var ErrAlreadyExists = errors.New("API already registered with different surname")

type Upstream struct {
	Target  string
	Surname string
}

type Janus struct {
	disabled bool
	upstream Upstream
	client   internal.Client
}

func NewNoOp() *Janus {
	return &Janus{
		disabled: true,
	}
}

func New(endpoint string, upstream Upstream) *Janus {
	return &Janus{
		upstream: upstream,
		client: internal.Client{
			endpoint,
		},
	}
}

// RegisterChi takes router and registers all endpoints in janus
func (j *Janus) RegisterChi(r chi.Router) error {
	if j.disabled {
		return nil
	}

	// walk the router without hitting janus
	services, err := internal.NewChi(r).Services()
	if err != nil {
		return errors.Wrap(err, "failed to walk chi router")
	}

	for _, candidate := range services {
		if err := j.register(candidate); err != nil {
			return errors.Wrap(err, "failed to register service")
		}
	}

	return nil
}

func (j *Janus) register(service internal.Service) error {
	// check if service already exists
	remote, err := j.client.GetAPI(service.Name)
	if err != nil {
		return errors.Wrap(err, "failed to get remote service")
	}
	if remote != nil {
		if remote.Surname != j.upstream.Surname {
			return ErrAlreadyExists
		}

		//  check if upstream is duplicate
		for _, target := range remote.Proxy.Upstreams.Targets {
			if target.Target == j.upstream.Target {
				return nil
			}
		}

		// modify remote service
		remote.Proxy.Upstreams.Targets = append(
			remote.Proxy.Upstreams.Targets,
			internal.Target{Target: j.upstream.Target})

		if err := j.client.UpdateAPI(remote.Name, remote); err != nil {
			return errors.Wrap(err, "failed to update remote service")
		}
		return nil
	}

	// add new service definition
	service.Surname = j.upstream.Surname
	service.Proxy.Upstreams = internal.Upstreams{
		Balancing: "weight",
		Targets:   []internal.Target{{Target: j.upstream.Target, Weight: 10}},
	}
	if err := j.client.AddAPI(&service); err != nil {
		return errors.Wrap(err, "failed to add new service")
	}
	return nil
}

func (j *Janus) RegisterGojiEndpoint(endpoint, method string) error {
	if j.disabled {
		return nil
	}
	route := internal.GetRouteForGoji(endpoint)
	name := internal.GetName(route, method)
	methods := []string{method}
	service := internal.Service{
		Name:   name,
		Active: true,
		Proxy: internal.Proxy{
			AppendPath: true,
			ListenPath: route,
			Methods:    methods,
		},
	}

	if err := j.register(service); err != nil {
		return errors.Wrap(err, "failed to register service")
	}
	return nil
}
