package main

import (
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

	log.Initialize(config.GetLogLevel())
	if err := data.Initialize(); err != nil {
		panic(err)
	}

	server := &http.Server{
		Addr: "0.0.0.0:8888",
		Handler: httpserver.Router(),
	}

	log.Info("Starting Server ...")
	defer log.Info("Shutdown Server.")
	// TODO: implement graceful shutdown in httpserver and defer it instead of the above call.

	go func() {
		if err := server.ListenAndServe(); err != nil{
			log.Error("Cannot start the http server.", err)
			os.Exit(1)
		}
	}()

	// wait for os.Interrupt
	s := <-c
	log.Info("Received signal:", s)
}
