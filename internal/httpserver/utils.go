package httpserver

import (
	"github.com/pokekrishna/chitchat/pkg/log"
	"net/http"
)

func writeErrorToClient(message string, w http.ResponseWriter){
	if _, err := w.Write([]byte(message)); err != nil {
		log.Error("cannot write error message to client:", message)
	}
}
