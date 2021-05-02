package main

import (
	"github.com/pokekrishna/chitchat/internal/config"
	"github.com/pokekrishna/chitchat/internal/data"
	"github.com/pokekrishna/chitchat/internal/httpserver"
	"github.com/pokekrishna/chitchat/pkg/log"
	"net/http"
)

func main() {
	log.Initialize(config.GetLogLevel())
	if err := data.Initialize(); err != nil {
		panic(err)
	}

	server := &http.Server{
		Addr: "0.0.0.0:8888",
		Handler: httpserver.Router(),
	}

	log.Info("Starting Server ...")
	server.ListenAndServe()

}
