package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/gorilla/handlers"

	"github.com/kiran-golang/jfrog-stats/api"
	"github.com/kiran-golang/jfrog-stats/config"
)

func main() {

	// Create router for the API URIs supported
	httpRouter := api.NewRouter()

	// Log all the http requests to stdout
	loggedRouter := handlers.LoggingHandler(os.Stdout, httpRouter)

	// Start service on configured port.
	httpServer := &http.Server{
		Handler: loggedRouter,
		Addr:    ":" + config.GetConfiguration().ServicePort,
	}

	// This section of the code handles graceful
	// shutdown of the service
	connectionsClose := make(chan struct{})
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		log.Println("Shutting Down Service")
		httpServer.Shutdown(context.Background())
		close(connectionsClose)
	}()

	tlsConfig, err := config.GetTLSConfig()
	if err != nil {
		log.Println("Unable to Load Certificates")
		log.Println("Starting Jfrog-stats Service as http service on " + config.GetConfiguration().ServicePort)
		log.Fatal(httpServer.ListenAndServe())
	} else {
		httpServer.TLSConfig = tlsConfig
		// empty strings because tlsconfig already has this information
		log.Println("Starting Jfrog-stats Service as https service on " + config.GetConfiguration().ServicePort)
		log.Fatal(httpServer.ListenAndServeTLS("", ""))
	}
}
