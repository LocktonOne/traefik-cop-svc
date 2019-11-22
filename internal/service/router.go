package service

import (
	"github.com/go-chi/chi"
	"gitlab.com/tokend/traefik-cop/internal/config"
	"net/http"
)

func Router(s *Service, cfg config.Config) chi.Router {
	r := chi.NewRouter()

	r.Use(
		ape.RecoverMiddleware(cfg.Log()),
		ape.LoganMiddleware(cfg.Log()),
		ape.CtxMiddleWare(),
	)

	r.Route("/integrations/traefik", func(r chi.Router) {
		r.Post("/services", s.addService)
	})

	return r
}

func (s *Service) addService(w http.ResponseWriter, r *http.Request) {

}