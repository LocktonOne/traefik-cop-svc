package service

import (
	"context"
	"gitlab.com/distributed_lab/kit/traefik"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/traefik-cop/internal/config"
	"net"
	"net/http"
)

type Service struct {
	config config.Config
	listener         net.Listener
	log *logan.Entry
	trefik *traefik.Traefik
}

func NewService(cfg config.Config) *Service {
	return &Service{
		config: cfg,
		log: cfg.Log(),
		trefik: cfg.Traefik(),
	}
}

func(s *Service) Run(ctx context.Context) {
	// service's API
	err := s.runService()
	if err != nil {
		panic(err)
	}
}

func(s *Service) runService() error {
	r := s.router()
	err := http.Serve(s.listener, r)
	if err != nil {
		return errors.Wrap(err, "server stopped with error")
	}

	return nil
}