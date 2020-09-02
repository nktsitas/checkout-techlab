package router

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"

	"github.com/nktsitas/checkout-techlab/handlers"
	"github.com/nktsitas/checkout-techlab/logger"
	"github.com/nktsitas/checkout-techlab/auth"

	_ "github.com/nktsitas/checkout-techlab/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

func NewRouter() *mux.Router {
	routes := CreateRoutes()

	// Create a new mux.Router
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/login", auth.Login).Methods("POST")
	router.HandleFunc("/status/ping", handlers.Ping).Methods("GET")

	router.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)

	// Add logging & authentication middleware and register routes
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = auth.Authenticate(handler, route.Name)
		handler = logger.APICallsLogger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

func CreateRoutes() []Route {
	var routes []Route

	routes = append(routes, Route{"CreateAuthorization", "POST", "/authorize", handlers.CreateAuthorizationHandler})
	routes = append(routes, Route{"Void", "POST", "/void", handlers.VoidHandler})
	routes = append(routes, Route{"Capture", "POST", "/capture", handlers.CaptureHandler})
	routes = append(routes, Route{"Refund", "POST", "/refund", handlers.RefundHandler})

	log.WithFields(log.Fields{
		"routes": routes,
	}).Debug("Routes initialized")

	return routes
}