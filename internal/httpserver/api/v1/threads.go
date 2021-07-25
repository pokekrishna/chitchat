package v1

import (
	"encoding/json"
	"github.com/pokekrishna/chitchat/internal/data"
	"github.com/pokekrishna/chitchat/pkg/content"
	"github.com/pokekrishna/chitchat/pkg/log"
	"net/http"
)

// TODO: does the package structure needs better naming like 'models' ?

// TODO: how to deal with multiple http methods on same model ...
// TODO: ... like -X GET Threads, and -X POST Threads.
func Threads(app *data.App, w http.ResponseWriter, r *http.Request) {
	var respBody []byte
	threads, err := app.Threads()
	if err != nil {
		log.Error("Cannot get threads", err)
	}
	contentType, _ := content.ExtractContentType(r)
	switch contentType {
	case content.TypeJSON, content.TypeNotSupported:
		var err error
		respBody, err =json.Marshal(threads)
		if err != nil {
			log.Error("Cannot marshal threads", err)
		}
	}

	// TODO: Is it a design flaw to simply dump to resp from db?
	_, err = w.Write(respBody)
	if err != nil {
		log.Error("Error writing to response", err)
	}
}