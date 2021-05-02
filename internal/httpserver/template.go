package httpserver

import (
	"fmt"
	"github.com/pokekrishna/chitchat/pkg/log"
	"text/template"
)

func parseTemplateFiles(filenames ... string) (t *template.Template){
	var files [] string
	t = template.New("layout")
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("internal/httpserver/templates/%s", file))
	}
	log.Info( "Parsing these files:", files)
	t = template.Must(template.ParseFiles(files ...))
	return
}