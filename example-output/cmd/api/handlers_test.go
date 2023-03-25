package main

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

