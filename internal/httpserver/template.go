package httpserver

import (
	"fmt"
	"github.com/pokekrishna/chitchat/pkg/log"
	"net/http"
	"text/template"
)

func parseTemplateFiles(filenames ...string) (*template.Template, error){
	var files [] string
	t := template.New("layout")
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("internal/httpserver/templates/%s", file))
	}
	log.Info( "Parsing these files:", files)
	t, err := template.ParseFiles(files ...)
	if err != nil{
		return nil, err
	}
	return t, nil
}

// generateHTML and parses templates and writes to http http.ResponseWriter
func generateHTML(w http.ResponseWriter, data interface{}, filenames ...string) error{
	t, err := parseTemplateFiles(filenames...)
	if err != nil{
		return err
	}
	err = t.Execute(w, data)
	if err != nil {
		return err
	}
	return nil
}