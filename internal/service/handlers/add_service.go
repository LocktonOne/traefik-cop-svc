package handlers

import (
	"net/http"
	"time"

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
						URL:      request.Data.Attributes.Url,
						Scheme:   "http",
						Port:     request.Data.Attributes.Port,
						Interval: durationWithFallback(request.Data.Attributes.HealthcheckInterval, TraefikCfg(r).GeneralHealthcheckInterval),
						Timeout:  durationWithFallback(request.Data.Attributes.HealthcheckTimout, TraefikCfg(r).GeneralHealthcheckTimeout),
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

func durationWithFallback(value *string, fallbackValue time.Duration) time.Duration {
	if value == nil {
		return fallbackValue
	}
	d, err := time.ParseDuration(*value)
	if err != nil {
		return fallbackValue
	}
	if d == 0 {
		return fallbackValue
	}
	return d
}
