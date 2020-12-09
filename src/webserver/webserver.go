package webserver

import (
	"context"
	"encoding/json"
	"fmt"
	"honeypot/settings"
	"log"
	"net/http"
	_ "net/http/pprof"
)

func index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	e, ok := r.Context().Value("env").(*env)
	if !ok {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}
	connAttemps, err := e.tl.GetConnAttemps(r.Context())
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, i := range connAttemps {
		fmt.Println(i)
	}

	// http.ServeFile(w, r, "cmd/webserver/index.html")
	js, err := json.Marshal(connAttemps)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func injectEnv(e *env, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "env", e)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func ServerMain() error {
	e := newEnv()
	defer e.destroy()

	log.Println("Webserver starting")
	http.HandleFunc("/", injectEnv(e, index))
	err := http.ListenAndServe(settings.WebserverAddr, nil)
	if err != nil {
		return err
	}
	return nil
}
