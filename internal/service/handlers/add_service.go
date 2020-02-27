package handlers

import (
	"net/http"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/tokend/traefik-cop/internal/service/requests"
	"gitlab.com/tokend/traefik-cop/internal/service/traefik"
	traefik2 "gitlab.com/tokend/traefik-cop/traefik"
)

func AddService(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewAddServiceRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	err = Updater(r, traefik.Backend{
		Router: traefik2.Router{
			Service: request.Data.ID,
			Rule:    request.Data.Attributes.Rule,
		},
		Service: traefik2.Service{
			LoadBalancer: traefik2.ServersLoadBalancer{
				Servers: []traefik2.Server{
					{
						URL:    request.Data.Attributes.Url,
						Scheme: "http",
						Port:   request.Data.Attributes.Port,
					},
				},
			},
		},
	})

	if err != nil {
		Log(r).WithError(err).Error("failed to register service")
		ape.RenderErr(w, problems.InternalError())
		return
	}
}
