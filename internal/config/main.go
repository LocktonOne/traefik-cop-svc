package config

import (
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/kit/traefik"
)

type Config interface {
	traefik.Traefiker
	comfig.Logger
}

type config struct {
	traefik.Traefiker
	comfig.Logger
	getter kv.Getter
}

func NewConfig(getter kv.Getter) Config {
	return &config{
		getter:     getter,
		Logger:     comfig.NewLogger(getter, comfig.LoggerOpts{Release: "unverified"}),
		Traefiker:  traefik.NewTraefiker(getter),
	}
}

