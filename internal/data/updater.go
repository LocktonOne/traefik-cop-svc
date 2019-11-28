package data

import "gitlab.com/tokend/traefik-cop/internal/service/traefik"

type Updater func(backend traefik.Backend) error
