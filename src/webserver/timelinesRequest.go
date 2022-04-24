package webserver

import (
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

func getBytes(w http.ResponseWriter, r *http.Request) {
	serviceValue, err := getKeyFromURL(r, "service")
	if err != nil {
		yieldInvalidService(w)
		return
	}

	ip, err := getIPFromServiceName(serviceValue)
	if err != nil {
		yieldInvalidService(w)
		return
	}

	queryer := func(e *env, rangeValue string) (interface{}, error) {
		return e.tl.GetBytes(r.Context(), rangeValue, ip)
	}
	timelinesRequest(w, r, queryer)
}

func exportData(w http.ResponseWriter, r *http.Request) {
	timelinesRequestCSV(w, r, "conn_attemps.csv", func(e *env, rangeValue string) (interface{}, error) {
		return e.tl.ExportData(r.Context())
	})
}

func getTimeRange(r *http.Request) string {
	rangeValue, err := getQueryParamRange(r)
	if err != nil || rangeValue == "" {
		// TODO: Handle this better
		return "mo"
	}
	return rangeValue
}

func timelinesRequest(
	w http.ResponseWriter,
	r *http.Request,
	queryer func(*env, string) (interface{}, error),
) {
	ret := timelinesRequestCommon(w, r, queryer)
	if ret == nil {
		return
	}
	responde(w, r, ret)
}

func timelinesRequestCSV(
	w http.ResponseWriter,
	r *http.Request,
	csvFilename string,
	queryer func(*env, string) (interface{}, error),
) {
	ret := timelinesRequestCommon(w, r, queryer)
	if ret == nil {
		return
	}
	respondeCSV(w, r, csvFilename, ret.(string))
}

func timelinesRequestCommon(
	w http.ResponseWriter,
	r *http.Request,
	queryer func(*env, string) (interface{}, error),
) interface{} {
	e, ok := r.Context().Value("env").(*env)
	if !ok {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return nil
	}

	rangeValue := getTimeRange(r)
	ret, err := queryer(e, rangeValue)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}

	return ret
}

func yieldInvalidService(w http.ResponseWriter) {
	http.Error(w, "Invalid service", http.StatusBadRequest)
}
