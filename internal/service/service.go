package service

import (
	"context"
	"net"
	"net/http"
	"sync"

	"gitlab.com/tokend/traefik-cop/internal/data"

	"gitlab.com/tokend/traefik-cop/traefik"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/traefik-cop/internal/config"
	traefik2 "gitlab.com/tokend/traefik-cop/internal/service/traefik"
)

type Service struct {
	config   config.Config
	listener net.Listener
	log      *logan.Entry
	traefik  *traefik.Traefik
	updater  data.Updater

	*sync.RWMutex
	backends map[string]traefik2.Backend
}

func NewService(cfg config.Config) *Service {
	return &Service{
		config:   cfg,
		RWMutex:  &sync.RWMutex{},
		log:      cfg.Log(),
		listener: cfg.Listener(),
		traefik:  cfg.Traefik(),
		backends: make(map[string]traefik2.Backend),
	}
}

func (s *Service) Run(ctx context.Context) {
	s.updater = func(backend traefik2.Backend) error {
		s.Lock()
		defer s.Unlock()

		s.updateConfiguration(backend)

		//TODO send conf request
		err := s.register()
		if err != nil {
			return errors.Wrap(err, "failed to register configuration")
		}
		return nil
	}

	// service's API
	err := s.runService()
	if err != nil {
		panic(err)
	}
}

func (s *Service) register() error {
	routers := make(map[string]traefik.Router)
	services := make(map[string]traefik.Service)

	for name, backend := range s.backends {
		routers[name] = backend.Router
		services[name] = backend.Service
	}

	err := s.traefik.RegisterConfiguration(
		traefik.Configuration{
			HTTP: traefik.HTTPConfiguration{
				Routers:  routers,
				Services: services,
			},
		})

	return err
}

func (s *Service) updateConfiguration(backend traefik2.Backend) {
	existing, ok := s.backends[backend.Router.Service]
	if !ok {
		s.backends[backend.Router.Service] = backend
		return
	}

	for _, server := range existing.Service.LoadBalancer.Servers {
		// server already registered
		if server == backend.Service.LoadBalancer.Servers[0] {
			return
		}
	}

	existing.Service.LoadBalancer.Servers = append(existing.Service.LoadBalancer.Servers, backend.Service.LoadBalancer.Servers[0])
	s.backends[backend.Router.Service] = existing
}

func (s *Service) runService() error {
	r := s.router()
	err := http.Serve(s.listener, r)
	if err != nil {
		return errors.Wrap(err, "server stopped with error")
	}

	return nil
}
