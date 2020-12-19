package webserver

import (
	"encoding/json"
	"net/http"
)

func responde(w http.ResponseWriter, r *http.Request, ret interface{}) {
	js, err := json.Marshal(ret)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
