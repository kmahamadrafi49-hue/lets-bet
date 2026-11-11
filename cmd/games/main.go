r := mux.NewRouter()

	// FORCE OVERRIDE: Serve index.html and static files robustly
	r.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		if _, err := os.Stat("index.html"); err == nil {
			http.ServeFile(w, req, "index.html")
		} else if _, err := os.Stat("cmd/games/index.html"); err == nil {
			http.ServeFile(w, req, "cmd/games/index.html")
		} else {
			http.ServeFile(w, req, "index.html")
		}
	})

	// Static files mariyu images kosam correct file server mapping
	r.PathPrefix("/images/").Handler(http.FileServer(http.Dir(".")))
	r.PathPrefix("/static/").Handler(http.FileServer(http.Dir(".")))
	r.PathPrefix("/row-").Handler(http.FileServer(http.Dir(".")))
