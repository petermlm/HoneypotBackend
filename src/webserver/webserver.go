package webserver

import (
	"context"
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

func getTotalConsumptions(w http.ResponseWriter, r *http.Request) {
	timelinesRequest(w, r, func(e *env) (interface{}, error) {
		return e.tl.GetTotalConsumptions(r.Context())
	})
}

func getMapData(w http.ResponseWriter, r *http.Request) {
	timelinesRequest(w, r, func(e *env) (interface{}, error) {
		return e.tl.GetMapData(r.Context())
	})
}

func getConnAttmps(w http.ResponseWriter, r *http.Request) {
	timelinesRequest(w, r, func(e *env) (interface{}, error) {
		return e.tl.GetConnAttemps(r.Context())
	})
}

func getTopConsumers(w http.ResponseWriter, r *http.Request) {
	timelinesRequest(w, r, func(e *env) (interface{}, error) {
		return e.tl.GetTopConsumers(r.Context())
	})
}

func getTopFlavours(w http.ResponseWriter, r *http.Request) {
	timelinesRequest(w, r, func(e *env) (interface{}, error) {
		return e.tl.GetTopFlavours(r.Context())
	})
}

func timelinesRequest(w http.ResponseWriter, r *http.Request, queryer func(*env) (interface{}, error)) {
	e, ok := r.Context().Value("env").(*env)
	if !ok {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}
	ret, err := queryer(e)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	responde(w, r, ret)
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

	router.HandleFunc("/totalConsumptions", injectEnv(e, getTotalConsumptions)).Methods("GET")
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
