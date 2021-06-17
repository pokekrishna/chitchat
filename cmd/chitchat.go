package main

import (
	"context"
	"github.com/pokekrishna/chitchat/internal/config"
	"github.com/pokekrishna/chitchat/internal/data"
	"github.com/pokekrishna/chitchat/internal/httpserver"
	"github.com/pokekrishna/chitchat/pkg/log"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	log.Initialize(config.LogLevel())
	db, err := data.Initialize()
	if err != nil {
		panic(err)
	}

	server := &http.Server{
		Addr:    "0.0.0.0:8888",
		Handler: httpserver.Router(db),
	}

	ctx := context.Background()

	log.Info("Starting Server ...")
	defer httpserver.Shutdown(ctx, server)

	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Error("Cannot start the http server.", err)
			os.Exit(1)
		}
	}()

	// wait for os.Interrupt
	s := <-c
	log.Info("Received signal:", s)
}
