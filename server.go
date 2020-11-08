package main

import (
	"log"
	"net/http"
	"os"

	"github.com/vulcand/oxy/forward"
	"github.com/vulcand/oxy/testutils"
)

func main() {
	url := getURL()
	log.Println("Routing requests to:", url)
	// Forwards incoming requests to whatever location URL points to, adds proper forwarding headers
	fwd, _ := forward.New()

	redirect := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// let us forward this request to another server
		req.URL = testutils.ParseURI(url)
		fwd.ServeHTTP(w, req)
	})

	// that's it! our reverse proxy is ready!
	s := &http.Server{
		Addr:    ":"+getPort(),
		Handler: redirect,
	}
	s.ListenAndServe()
}

func getURL() string {
	metricsServerURL := os.Getenv("METRICS_SERVER_URL")

	if metricsServerURL == "" {
		panic("Metrics server url not found in environment variables. Please set METRICS_SERVER_URL env variable.")
	}

	return metricsServerURL
}

func getPort() string {
	port, exists := os.LookupEnv("PORT")
	if !exists || port == "" {
		port = "8080"
	}
	return port
}
