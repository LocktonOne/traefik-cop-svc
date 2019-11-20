package comfig

import (
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/kit/janus"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type Januser interface {
	Janus() *janus.Janus
}

func NewJanuser(getter kv.Getter) Januser {
	return &januser{
		getter: getter,
	}
}

type januser struct {
	getter kv.Getter
	once   Once
}

func (j *januser) Janus() *janus.Janus {
	return j.once.Do(func() interface{} {
		raw := kv.MustGetStringMap(j.getter, "janus")

		var probe struct {
			Disabled bool `fig:"disabled"`
		}

		if err := figure.Out(&probe).From(raw).Please(); err != nil {
			panic(errors.Wrap(err, "failed to figure out janus probe"))
		}

		if probe.Disabled {
			return janus.NewNoOp()
		}

		var config struct {
			Endpoint string `fig:"endpoint,required"`
			Upstream string `fig:"upstream,required"`
			Surname  string `fig:"surname,required"`
		}

		if err := figure.Out(&config).From(raw).Please(); err != nil {
			panic(errors.Wrap(err, "failed to figure out janus"))
		}

		return janus.New(config.Endpoint, janus.Upstream{
			Target:  config.Upstream,
			Surname: config.Surname,
		})
	}).(*janus.Janus)
}
