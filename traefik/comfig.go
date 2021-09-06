package traefik

import (
	"context"
	"reflect"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/cast"

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
			RestDisabled  bool `fig:"rest_disabled"`
			RedisDisabled bool `fig:"redis_disabled"`
		}

		if err := figure.Out(&probe).From(raw).Please(); err != nil {
			panic(errors.Wrap(err, "failed to figure out traefik probe"))
		}

		if probe.RestDisabled && probe.RedisDisabled {
			return NewNoOp()
		}

		var config TraefikConfig
		fig := figure.Out(&config).From(raw)
		if !probe.RedisDisabled { // doing this to keep configs 'disabled'-compatible
			fig = fig.With(figure.BaseHooks, redisClientHooks)
		}

		if err := fig.Please(); err != nil {
			return errors.Wrap(err, "failed to figure out traefik")
		}

		switch {
		case probe.RedisDisabled:
			return NewWithRestInit(config)
		case probe.RestDisabled:
			return NewWithRedisInit(config)
		}

		return NewNoOp()
	}).(*Traefik)
}

var redisClientHooks = figure.Hooks{
	"*redis.Client": func(value interface{}) (reflect.Value, error) {
		raw, err := cast.ToStringMapE(value)
		if err != nil {
			return reflect.Value{}, errors.Wrap(err, "failed to parse map[string]interface{}")
		}

		config := struct {
			Addr     string `fig:"address"`
			Password string `fig:"password"`
			DB       int    `fig:"db"`
		}{
			Addr:     "127.0.0.1:6379",
			Password: "",
			DB:       0,
		}

		if err := figure.Out(&config).From(raw).Please(); err != nil {
			panic(errors.Wrap(err, "failed to get data redis from config"))
		}

		clientRedis := redis.NewClient(&redis.Options{
			Addr:     config.Addr,
			Password: config.Password,
			DB:       config.DB,
		})

		if err := clientRedis.Ping(context.TODO()).Err(); err != nil {
			panic(errors.Wrap(err, "failed to connect to redis"))
		}

		return reflect.ValueOf(clientRedis), nil

	},
}
