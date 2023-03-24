package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/julienschmidt/httprouter"
)

// Struct encapsulating all configuration settings for application
type config struct {
	port int // Network port that we want server to listen on
}

// Struct encapsulating all dependancies for HTTP handlers, helpers and middleware
// Should also contain any variables pertaining to application state that needs
// to be accessible to any handlers
type application struct {
	config config // Configuration settings for application
}

func main() {
	var cfg config // Application configuration settings

	// Parse port and operating environment from given flags
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.Parse()

	// Logger to control messaging to stdout stream
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	app := &application{
		config: cfg,
	}

	// HTTP server with basic sensible timeout settings
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// Start HTTP Server
	logger.Printf("starting server on %s", srv.Addr)
	err := srv.ListenAndServe()
	logger.Fatal(err)
}

func (a *application) routes() *httprouter.Router {
	// Create a new HTTP router
	router := httprouter.New()
	// Attach endpoint handler methods
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", a.healthcheckHandler)
	// Return the configured router
	return router
}
