package main

import (
	"net/http"
)

func createHandler() http.Handler {

	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir(set.Root)))

	/*mux.HandleFunc("/ip",func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte(r.RemoteAddr))
	})*/

	return mux
}
