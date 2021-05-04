package httpserver

import "net/http"

func Router() *http.ServeMux {
	mux := http.NewServeMux()

	// The following directory is relative to the location where the
	// main program is being run from.
	files := http.FileServer(http.Dir("internal/httpserver/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", files))


	mux.HandleFunc("/", index)
	mux.HandleFunc("/err", errHandler)
	mux.HandleFunc("/login", login)
	mux.HandleFunc("/logout", logout)
	mux.HandleFunc("/authenticate", authenticate)

	return mux
}
