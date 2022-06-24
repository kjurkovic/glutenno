package router

import (
	"auth/config"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Router struct{}

type RouteHandler interface {
	pathPrefix() string
	setup(r *mux.Router)
}

func (r *Router) Setup(log *log.Logger, config *config.Config, router *mux.Router) {
	routeHandlers := []RouteHandler{
		&AuthRouter{log, config},
		&UserRouter{log, config},
	}

	for _, handler := range routeHandlers {
		if len(handler.pathPrefix()) == 0 {
			handler.setup(router)
		} else {
			handler.setup(router.PathPrefix(handler.pathPrefix()).Subrouter())
		}
	}
}

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Access-Control-Allow-Origin", "*")
		rw.Header().Set("Access-Control-Allow-Credentials", "true")
		rw.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		rw.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, PUT, OPTIONS")

		if r.Method == http.MethodOptions {
			rw.WriteHeader(204)
			return
		}

		next.ServeHTTP(rw, r)
	})
}
