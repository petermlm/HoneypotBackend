package webserver

import (
	"encoding/json"
	"fmt"
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

func respondeCSV(w http.ResponseWriter, r *http.Request, csvFilename, csvData string) {
	headerValue := fmt.Sprintf(`attachment; filename="%s"`, csvFilename)
	w.Header().Set("Content-Disposition", headerValue)
	w.Write([]byte(csvData))
}
