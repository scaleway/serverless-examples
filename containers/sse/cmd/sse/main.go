package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/r3labs/sse/v2"
)

var (
	// Serverless Containers will scale to multiple instances.
	// This is a way to identify which instance is sending the message to identify the source.
	hostName, _ = os.Hostname()
)

func main() {
	e := echo.New()

	server := sse.New()             // create SSE broadcaster server
	server.AutoReplay = false       // do not replay messages for each new subscriber that connects
	_ = server.CreateStream("time") // EventSource in "index.html" connecting to stream named "time"

	go func(s *sse.Server) {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				s.Publish("time", &sse.Event{
					Data: []byte("time: " + time.Now().Format(time.RFC3339Nano) + " from: " + hostName),
				})
			}
		}
	}(server)

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Static("/", "www")

	e.GET("/sse", func(c echo.Context) error { // longer variant with disconnect logic
		log.Printf("The client is connected: %v\n", c.RealIP())
		go func() {
			<-c.Request().Context().Done() // Received Browser Disconnection
			log.Printf("The client is disconnected: %v\n", c.RealIP())
		}()

		server.ServeHTTP(c.Response(), c.Request())

		return nil
	})

	if err := e.Start(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
