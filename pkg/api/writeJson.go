package api

import (
	"encoding/json"
	"net/http"
)

func writeJson(res http.ResponseWriter, data any) {
	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(data)
}

func writeJsonError(res http.ResponseWriter, err error, code int) {
	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	res.WriteHeader(code)
	json.NewEncoder(res).Encode(map[string]string{"error": err.Error()})
}
