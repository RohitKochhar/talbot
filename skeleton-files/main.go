package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// Application version number
const version = "1.0.0"

// Struct encapsulating all configuration settings for application
type config struct {
	port int    // Network port that we want server to listen on
	env  string // Current operating environment
}

// Struct encapsulating all dependancies for HTTP handlers, helpers and middleware
type application struct {
	config config      // Configuration settings for application
	logger *log.Logger // Logger to control messaging
}

func main() {
	var cfg config // Application configuration settings

	// Parse port and operating environment from given flags
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.Parse()

	// Logger to control messaging to stdout stream
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	app := &application{
		config: cfg,
		logger: logger,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/healthcheck", app.healthcheckHandler)

	// HTTP server with basic sensible timeout settings
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      mux,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// Start HTTP Server
	logger.Printf("starting %s server on %s", cfg.env, srv.Addr)
	err := srv.ListenAndServe()
	logger.Fatal(err)
}
