package traefik

import "gitlab.com/tokend/traefik-cop/traefik"

type Backend struct {
	Router  traefik.Router
	Service traefik.Service
}
