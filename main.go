package main

import (
	"fmt"
	"net/http"

	"github.com/nktsitas/checkout-techlab/router"
	"github.com/nktsitas/checkout-techlab/db"
	"github.com/nktsitas/checkout-techlab/gateway"
	
	log "github.com/sirupsen/logrus"
)

const port = "2012"

// @title Checkout.com API Challenge
// @version 1.0
// @description This is a simple Gateway service for Payments
// @termsOfService http://swagger.io/terms/
// @contact.name Nikos Tsitas
// @contact.email nktsitas@gmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:2012
// @BasePath /
func main() {
	db.DB = db.InitMemoryDB()
	gateway.Gateway = new(gateway.GatewayS)

	router := router.NewRouter()

	// Fire up server
	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	log.Info(fmt.Sprintf("Checkout Tech Test API - Candidate: Nikos Tsitas"))
	log.Info(fmt.Sprintf("Checkout Tech Test API - Listening on port: %s", port))

	if err := server.ListenAndServe(); err != nil {
		log.WithField("err", err).Fatal("Error initializing server")
	}
}