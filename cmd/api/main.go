package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/nhudson/longbox/internal/api"
	"github.com/oklog/run"
	"github.com/sirupsen/logrus"
)

const (
	defaultHTTPPort = 7575
	defaultURL      = "https://getcomics.info"
)

var (
	httpPort  = flag.Int("http-port", getEnvOrFallbackInt("PORT", defaultHTTPPort), "HTTP API Port")
	comicsURL = flag.String("comics-url", defaultURL, "Default URL to use right now only use https://getcomics.info")
)

func getEnvOrFallbackInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if value, err := strconv.Atoi(value); err == nil {
			return value
		}
	}
	return fallback
}

func main() {
	log := logrus.StandardLogger()

	var g run.Group

	{
		done := make(chan struct{})
		listenAddr := fmt.Sprintf(":%d", *httpPort)

		srv := &http.Server{
			Addr:         listenAddr,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
			IdleTimeout:  120 * time.Second,
			Handler:      api.NewServer(log.WithField("context", "api"), *comicsURL),
		}

		g.Add(func() error {
			log.Println("server is ready to handle requests at", listenAddr)
			err := srv.ListenAndServe()
			<-done
			log.Println("server stopped")
			return err
		}, func(error) {
			log.Println("server is shutting down...")

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			srv.SetKeepAlivesEnabled(false)
			if err := srv.Shutdown(ctx); err != nil {
				log.Fatalf("could not gracefully shutdown the server: %v", err)
			}
			close(done)
		})
	}

	if err := g.Run(); err != nil {
		log.Fatalln(err)
	}
}
