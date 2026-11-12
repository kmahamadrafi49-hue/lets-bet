package main

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	// Route for serving index.html correctly from /app directory
	r.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		if _, err := os.Stat("index.html"); err == nil {
			http.ServeFile(w, req, "index.html")
		} else if _, err := os.Stat("cmd/games/index.html"); err == nil {
			http.ServeFile(w, req, "cmd/games/index.html")
		} else {
			http.ServeFile(w, req, "index.html")
		}
	})

	// Static Files and Images mapping
	r.PathPrefix("/images/").Handler(http.FileServer(http.Dir(".")))
	r.PathPrefix("/static/").Handler(http.FileServer(http.Dir(".")))
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(".")))

	http.ListenAndServe(":8080", r)
}
