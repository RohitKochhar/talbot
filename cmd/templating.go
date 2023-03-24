package cmd

var MAIN_BASE = `package main

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
`

var HANDERS_BASE = `package main

import (
	"net/http"
)

// Replies with a 200 status indicating server is up and running
func (a *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	replyTextContent(w, r, http.StatusOK, "OK")
}

// replyTextContent wraps text content in a HTTP response and sends it
func replyTextContent(w http.ResponseWriter, r *http.Request, status int, content string) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(status)
	w.Write([]byte(content + "\n"))
}
`

var HANDLERS_TEST_BASE = `package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// setupAPI is a helper function that sets up
// the API for the tests, providing a cleanup function too
func setupAPI(t *testing.T) (string, func()) {
	t.Helper() // Mark the function as test helper
	app := &application{
		config: config{},
	}
	ts := httptest.NewServer(app.routes())
	return ts.URL, func() {
		ts.Close()
	}
}

// getHelper wraps the Get function in additional logic to
// assist with testing ease and clarity
func getHelper(t *testing.T, getUrl string, expBody string, expCode int) (r *http.Response) {
	r, err := http.Get(getUrl)
	if err != nil {
		t.Fatalf("error while sending GET request: %q", err)
	}
	// Check if the return code is what we expected
	if r.StatusCode != expCode {
		t.Fatalf("Expected %q, got %q.", http.StatusText(expCode),
			http.StatusText(r.StatusCode))
	}
	defer r.Body.Close()
	// We might not be expecting content
	if expBody != "" || expCode == http.StatusNotFound {
		// Check that we have the content that we expected
		var body []byte
		if body, err = io.ReadAll(r.Body); err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(body), expBody) {
			t.Fatalf("Expected %q, got %q.", expBody, string(body))
		}
	}

	return r
}

// TestIntegration runs sequential requests to the server to check the
// responses are what we would expect in real life usage
func TestIntegration(t *testing.T) {
	// Create a test server with routing defined in ./main.go
	url, cleanup := setupAPI(t)
	// Close server when testing is complete
	defer cleanup()
	// Check that the server is healthy
	_ = getHelper(t, url+"/v1/healthcheck", "OK", http.StatusOK)
}
`

// Creates main.go, handlers.go and handlers_test.go
func GenerateGoSourceFiles(target string) error {
	// Create main.go
	mainFile, err := createFile("cmd/api/main.go", target)
	if err != nil {
		return err
	}
	// Write main.go content to newly created file
	// ToDo: Add templating here where necessary
	if err := writeFile(mainFile, MAIN_BASE); err != nil {
		return err
	}
	// Create handlers.go
	handlersFile, err := createFile("cmd/api/handlers.go", target)
	if err != nil {
		return err
	}
	// Write handlers.go content to newly created file
	// ToDo: Add templating here where necessary
	if err := writeFile(handlersFile, HANDERS_BASE); err != nil {
		return err
	}
	// Create handlers_test.go
	handlersTestFile, err := createFile("cmd/api/handlers_test.go", target)
	if err != nil {
		return err
	}
	// Write handlers_test.go content to newly created file
	// ToDo: Add templating here where necessary
	if err := writeFile(handlersTestFile, HANDLERS_TEST_BASE); err != nil {
		return err
	}
	return nil
}
