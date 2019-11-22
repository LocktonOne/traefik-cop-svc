package internal

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"
)

type Chi struct {
	r chi.Router
}

func NewChi(r chi.Router) *Chi {
	return &Chi{r}
}

func (c *Chi) Routes() (string, error) {
	var routes []string
	walk := func(method, route string, _ http.Handler, _ ...func(http.Handler) http.Handler) error {
		route = strings.Replace(route, "/*/", "/", -1)
		if len(route) > 1 {
			route = strings.TrimRight(route, "/")
		}

		routes = append(routes, fmt.Sprintf("Path(`%s`)", route))
		return nil
	}

	err := chi.Walk(c.r, walk)
	if err != nil {
		return "", errors.Wrap(err, "failed to walk router")
	}
	return strings.Join(routes, "||"), nil
}
