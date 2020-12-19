package webserver

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func getKeyFromURL(r *http.Request, key string) (string, error) {
	vars := mux.Vars(r)
	if value, ok := vars[key]; ok {
		return value, nil
	}

	errStr := fmt.Sprintf("Not in request: %s", key)
	return "", errors.New(errStr)
}
