package app

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

type Route struct {
	Name       string
	Method     string
	Path       string
	Secure     bool
	HandleFunc http.HandlerFunc
}

type Routes []Route

func (a *App) configureRouter() {
	a.router.StrictSlash(true)
	a.router.Use(a.loggingMiddleware)

	metricRoutes := Routes{
		Route{
			Name:       "HeartBeat",
			Method:     "GET",
			Path:       "/heartbeat",
			Secure:     false,
			HandleFunc: a.heartbeat,
		},
	}

	for _, route := range metricRoutes {
		handler := route.HandleFunc

		a.router.
			Methods(route.Method).
			Path(route.Path).
			Name(route.Name).
			Handler(handler)
	}

	ordersRoutes := Routes{
		Route{
			Name:       "NewOrder",
			Method:     "POST",
			Path:       "/orders",
			Secure:     false,
			HandleFunc: a.newOrder,
		},
		Route{
			Name:       "AddItems",
			Method:     "POST",
			Path:       "/orders/{order_id:[a-z]{3,16}}/items",
			Secure:     false,
			HandleFunc: a.addItems,
		},
		Route{
			Name:       "GetOrder",
			Method:     "GET",
			Path:       "/orders/{order_id:[a-z]{3,16}}",
			Secure:     false,
			HandleFunc: a.getOrder,
		},
		Route{
			Name:       "DoneOrder",
			Method:     "POST",
			Path:       "/orders/{order_id:[a-z]{3,16}}/done",
			Secure:     false,
			HandleFunc: a.doneOrder,
		},
		Route{
			Name:       "GetALLOrders",
			Method:     "GET",
			Path:       "/orders/",
			Secure:     false,
			HandleFunc: a.getALLOrders,
		},
	}

	for _, route := range ordersRoutes {
		handler := route.HandleFunc

		a.router.
			Methods(route.Method).
			Path(route.Path).
			Name(route.Name).
			Handler(handler)
	}
}

func (a *App) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.logger.Infof("started %s %s", r.Method, r.RequestURI)

		wrapped := WrapResponseWriter(w)

		next.ServeHTTP(wrapped, r)

		var level logrus.Level
		switch {
		case wrapped.Status() >= 500:
			level = logrus.ErrorLevel
		case wrapped.Status() >= 400:
			level = logrus.WarnLevel
		default:
			level = logrus.InfoLevel
		}

		a.logger.Logf(
			level,
			"completed %s %s with %d %s",
			r.Method,
			r.RequestURI,
			wrapped.Status(),
			http.StatusText(wrapped.Status()),
		)
	})
}
