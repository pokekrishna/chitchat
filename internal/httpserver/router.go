package httpserver

import "net/http"

func Router() *http.ServeMux {
	mux := http.NewServeMux()

	// The following directory is relative to the location where the
	// main program is being run from.
	files := http.FileServer(http.Dir("internal/httpserver/static"))
	staticHandler:= http.StripPrefix("/static/", files)
	mux.HandleFunc("/static/", logHandler(staticHandler.(http.HandlerFunc)))


	mux.HandleFunc("/", logHandler(index))
	mux.HandleFunc("/err", logHandler(errHandler))
	mux.HandleFunc("/login", logHandler(login))
	mux.HandleFunc("/logout", logHandler(logout))
	mux.HandleFunc("/authenticate", logHandler(authenticate))

	return mux
}
