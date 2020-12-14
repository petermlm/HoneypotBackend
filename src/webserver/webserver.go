package webserver

import (
	"context"
	"encoding/json"
	"honeypot/settings"
	"log"
	"net/http"
	_ "net/http/pprof"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "cmd/webserver/index.html")
}

func getMapData(w http.ResponseWriter, r *http.Request) {
	e, ok := r.Context().Value("env").(*env)
	if !ok {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}
	mapData, err := e.tl.GetMapData(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	js, err := json.Marshal(mapData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func getConnAttmps(w http.ResponseWriter, r *http.Request) {
	e, ok := r.Context().Value("env").(*env)
	if !ok {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}
	connAttemps, err := e.tl.GetConnAttemps(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	js, err := json.Marshal(connAttemps)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func getTopConsumers(w http.ResponseWriter, r *http.Request) {
	e, ok := r.Context().Value("env").(*env)
	if !ok {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}
	mapData, err := e.tl.GetTopConsumers(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	js, err := json.Marshal(mapData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func getTopFlavours(w http.ResponseWriter, r *http.Request) {
	e, ok := r.Context().Value("env").(*env)
	if !ok {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}
	mapData, err := e.tl.GetTopFlavours(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	js, err := json.Marshal(mapData)
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
	log.Println("Webserver starting")

	e := newEnv()
	defer e.destroy()

	router := mux.NewRouter()
	router.HandleFunc("/", index).Methods("GET")
	router.HandleFunc("/map", injectEnv(e, getMapData)).Methods("GET")
	router.HandleFunc("/connAttemps", injectEnv(e, getConnAttmps)).Methods("GET")
	router.HandleFunc("/topConsumers", injectEnv(e, getTopConsumers)).Methods("GET")
	router.HandleFunc("/topFlavours", injectEnv(e, getTopFlavours)).Methods("GET")

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET"})

	err := http.ListenAndServe(settings.WebserverAddr, handlers.CORS(headersOk, originsOk, methodsOk)(router))
	if err != nil {
		return err
	}
	return nil
}
