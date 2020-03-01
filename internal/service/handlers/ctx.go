package handlers

import (
	"context"
	"net/http"

	traefik2 "gitlab.com/tokend/traefik-cop/traefik"

	"gitlab.com/tokend/traefik-cop/internal/data"

	"gitlab.com/tokend/traefik-cop/internal/service/traefik"

	"gitlab.com/distributed_lab/logan/v3"
)

type ctxKey int

const (
	logCtxKey ctxKey = iota
	updaterCtxkey
	traefikCfgCtxKey
)

func CtxLog(entry *logan.Entry) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, logCtxKey, entry)
	}
}

func Log(r *http.Request) *logan.Entry {
	return r.Context().Value(logCtxKey).(*logan.Entry)
}

func CtxUpdater(updater data.Updater) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, updaterCtxkey, updater)
	}
}

func Updater(r *http.Request, backend traefik.Backend) error {
	updater := r.Context().Value(updaterCtxkey).(data.Updater)
	return updater(backend)
}

func CtxTraefikCfg(c *traefik2.TraefikConfig) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, traefikCfgCtxKey, c)
	}
}

func TraefikCfg(r *http.Request) *traefik2.TraefikConfig {
	return r.Context().Value(traefikCfgCtxKey).(*traefik2.TraefikConfig)
}
