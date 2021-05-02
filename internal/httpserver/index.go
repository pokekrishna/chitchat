package httpserver

import (
	"fmt"
	"github.com/pokekrishna/chitchat/internal/data"
	"github.com/pokekrishna/chitchat/pkg/log"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request){
	var msg string
	threads, err := data.Threads()
	if err != nil {
		log.Error("Cannot get threads", err)
	}
	msg = fmt.Sprintf("%#v", threads)
	if _,err := w.Write([]byte(msg)); err != nil {
		log.Error("Error writing data to response writer", err)
	}
}
