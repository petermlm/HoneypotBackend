package webserver

import (
	"honeypot/timelines"
	"net/http"
)

func getTotalConsumptions(w http.ResponseWriter, r *http.Request) {
	timelinesRequest(w, r, func(e *env, rangeValue string) (interface{}, error) {
		return e.tl.GetTotalConsumptions(r.Context(), rangeValue)
	})
}

func getMapData(w http.ResponseWriter, r *http.Request) {
	timelinesRequest(w, r, func(e *env, rangeValue string) (interface{}, error) {
		return e.tl.GetMapData(r.Context(), rangeValue)
	})
}

func getConnAttmps(w http.ResponseWriter, r *http.Request) {
	timelinesRequest(w, r, func(e *env, rangeValue string) (interface{}, error) {
		return e.tl.GetConnAttemps(r.Context(), rangeValue)
	})
}

func getTopConsumers(w http.ResponseWriter, r *http.Request) {
	timelinesRequest(w, r, func(e *env, rangeValue string) (interface{}, error) {
		return e.tl.GetTopConsumers(r.Context(), rangeValue)
	})
}

func getTopFlavours(w http.ResponseWriter, r *http.Request) {
	timelinesRequest(w, r, func(e *env, rangeValue string) (interface{}, error) {
		return e.tl.GetTopFlavours(r.Context(), rangeValue)
	})
}

func getTimeRange(r *http.Request) string {
	rangeValue, err := getKeyFromURL(r, "range")
	if err != nil {
		return timelines.DefaultValidRanges
	}
	return rangeValue
}

func timelinesRequest(
	w http.ResponseWriter,
	r *http.Request,
	queryer func(*env, string) (interface{}, error),
) {
	e, ok := r.Context().Value("env").(*env)
	if !ok {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	rangeValue := getTimeRange(r)
	ret, err := queryer(e, rangeValue)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	responde(w, r, ret)
}
