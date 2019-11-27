package service

import (
	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/tokend/traefik-cop/internal/service/handlers"
)

func (s *Service) router() chi.Router {
	r := chi.NewRouter()

	r.Use(
		ape.RecoverMiddleware(s.config.Log()),
		ape.LoganMiddleware(s.config.Log()),
		ape.CtxMiddleware(
			handlers.CtxLog(s.log),
			handlers.CtxUpdater(s.updater)),
	)

	r.Route("/integrations/traefik", func(r chi.Router) {
		r.Post("/services", handlers.AddService)
	})

	return r
}
