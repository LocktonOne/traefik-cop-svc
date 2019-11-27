package traefik

import "gitlab.com/distributed_lab/kit/traefik"

type Backend struct {
	Router  traefik.Router
	Service traefik.Service
}
