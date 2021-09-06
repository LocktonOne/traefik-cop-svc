package handlers

import (
	"fmt"
	"net/http"

	"github.com/spf13/cast"
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

	fmt.Println("adding service", request.Data.Attributes.Name)
	fmt.Println("rule for service is", request.Data.Attributes.Rule)

	err = Updater(r, traefik.Backend{
		Router: traefik2.Router{
			Service:  request.Data.ID,
			Rule:     request.Data.Attributes.Rule,
			Priority: safeInt(request.Data.Attributes.RulePriority),
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

func safeInt(iptr *int32) int {
	if iptr == nil {
		return 0
	}
	return cast.ToInt(*iptr)
}
