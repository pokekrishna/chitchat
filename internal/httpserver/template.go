package httpserver

import (
	"fmt"
	"github.com/pokekrishna/chitchat/pkg/log"
	"net/http"
	"text/template"
)

func parseTemplateFiles(filenames ...string) (t *template.Template){
	var files [] string
	t = template.New("layout")
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("internal/httpserver/templates/%s", file))
	}
	log.Info( "Parsing these files:", files)
	t = template.Must(template.ParseFiles(files ...))
	return
}

// Genereates HTML and write to http responsewriter
// takes in writer, data for template, and files to join
func generateHTML(w http.ResponseWriter, data interface{}, filenames ...string) (err error){
	t := parseTemplateFiles(filenames...)
	err = t.Execute(w, data)
	if err != nil {
		return
	}
	return nil
}