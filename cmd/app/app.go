package main

import (
	"context"
	"flag"
	"github.com/timucingelici/news-gateway/internal/app/handlers"
	"github.com/timucingelici/news-gateway/pkg/config"
	"github.com/timucingelici/news-gateway/pkg/store"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {

	var wait time.Duration

	flag.DurationVar(
		&wait,
		"graceful-timeout",
		time.Second*15,
		"the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")

	flag.Parse()

	// Get config values
	conf, err := config.New()

	if err != nil {
		log.Fatalln("Failed to parse required env vars. Err : ", err)
	}

	// Setup the data store connection
	s, err := store.New(conf.RedisProtocol, conf.RedisAddr, conf.RedisReadTimeout, conf.RedisWriteTimeout)

	if err != nil {
		log.Fatalf("Failed to connect to the data store: %s\n", err)
	}

	// Setup the routes and the handlers
	handler := handlers.SetupRoutes(s)

	srv := &http.Server{
		Addr:         ":8090",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      handler,
	}

	go func() {
		log.Println("Starting the API server")
		if err := srv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Fatal("Failed to start API server")
			}
		}
	}()

	// Listen for an interrupt signal and try to shut down gracefully.
	shutdown(srv, wait)
}

func shutdown(srv *http.Server, wait time.Duration) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

	log.Println("Gracefully shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	err := srv.Shutdown(ctx)
	if err != nil {
		log.Fatal("Graceful shutdown has failed. Err: ", err)
	}

	log.Println("Shutdown complete.")
	os.Exit(0)
}
