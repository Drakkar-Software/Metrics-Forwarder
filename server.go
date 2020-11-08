package main

import (
	"log"
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
	log.Println("Routing requests to: ", url1)
	if err != nil {
		e.Logger.Fatal(err)
	}
	targets := []*middleware.ProxyTarget{
		{
			URL: url1,
		},
	}
	e.Use(middleware.Proxy(middleware.NewRoundRobinBalancer(targets)))

	e.Logger.Fatal(e.Start(":" + getPort()))
}

func getURL() string {
	metricsServerURL := os.Getenv("METRICS_SERVER_URL")

	if metricsServerURL == "" {
		panic("Metrics server url not found in environment variables. Please set METRICS_SERVER_URL env variable.")
	}
	if string(metricsServerURL[len(metricsServerURL)-1]) == "/" {
		metricsServerURL = metricsServerURL[:len(metricsServerURL)-1]
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

