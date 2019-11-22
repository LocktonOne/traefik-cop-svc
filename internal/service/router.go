package service

import (
	"github.com/go-chi/chi"
	"net/http"
	"gitlab.com/distributed_lab/ape"
)

func (s *Service) router() chi.Router {
	r := chi.NewRouter()

	r.Use(
		ape.RecoverMiddleware(s.config.Log()),
		ape.LoganMiddleware(s.config.Log()),
		ape.CtxMiddleWare(),
	)

	r.Route("/integrations/traefik", func(r chi.Router) {
		r.Post("/services", s.addService)
	})

	return r
}

func (s *Service) addService(w http.ResponseWriter, r *http.Request) {

}