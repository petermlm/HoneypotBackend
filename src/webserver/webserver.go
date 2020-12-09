package webserver

import (
	"honeypot/settings"
	"log"
	"net/http"
	_ "net/http/pprof"
)

func index(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "cmd/webserver/index.html")
}

func ServerMain() error {
	http.HandleFunc("/", index)
	err := http.ListenAndServe(settings.WebserverAddr, nil)
	if err != nil {
		return err
	}
	return nil
}
