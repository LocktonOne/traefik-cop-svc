package config

import (
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
)

type Config interface {
	comfig.Logger
}

type config struct {
	comfig.Logger
	getter kv.Getter
}

func NewConfig(getter kv.Getter) Config {
	return &config{
		getter:     getter,
		Logger:     comfig.NewLogger(getter),
	}
}

