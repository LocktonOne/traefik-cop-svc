package traefik

import (
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type Traefiker interface {
	Traefik() *Traefik
}

func NewTraefiker(getter kv.Getter) Traefiker {
	return &traefiker{
		getter: getter,
	}
}

type traefiker struct {
	getter kv.Getter
	once   comfig.Once
}

func (j *traefiker) Traefik() *Traefik {
	return j.once.Do(func() interface{} {
		raw := kv.MustGetStringMap(j.getter, "traefik")

		var probe struct {
			Disabled bool `fig:"disabled"`
		}

		if err := figure.Out(&probe).From(raw).Please(); err != nil {
			panic(errors.Wrap(err, "failed to figure out traefik probe"))
		}

		if probe.Disabled {
			return NewNoOp()
		}

		var config TraefikConfig

		if err := figure.Out(&config).From(raw).Please(); err != nil {
			panic(errors.Wrap(err, "failed to figure out janus"))
		}

		return New(config)
	}).(*Traefik)
}
