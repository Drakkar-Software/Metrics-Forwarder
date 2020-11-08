package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	//controller.Init()
	e := echo.New()
	e.Use(middleware.Logger())
	// Setup proxy
	url1, err := url.Parse(getURL())
	log.Println("Routing requests to:", url1)
	if err != nil {
		e.Logger.Fatal(err)
	}
	targets := []*middleware.ProxyTarget{
		{
			URL: url1,
		},
	}
	e.Use(getProxy(targets))

	e.Logger.Fatal(e.Start(":" + getPort()))
}

func getProxy(targets []*middleware.ProxyTarget) echo.MiddlewareFunc {
	c := middleware.DefaultProxyConfig

	// register forward urls
	c.Balancer = middleware.NewRoundRobinBalancer(targets)

	// force HTTP/1.1 protocol (from https://github.com/golang/go/issues/39302#issuecomment-635810949)
	tr := http.DefaultTransport.(*http.Transport).Clone()
	tr.ForceAttemptHTTP2 = false
	tr.TLSNextProto = make(map[string]func(authority string, c *tls.Conn) http.RoundTripper)
	tr.TLSClientConfig = &tls.Config{}
	c.Transport = tr

	return middleware.ProxyWithConfig(c)
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

