package api

import (
	"encoding/json"
	"net/http"
)

func writeJson(res http.ResponseWriter, data any) {
	resp, err := json.Marshal(data)
	if err != nil {
		http.Error(res, "Ошибка кодирования ответа в json", http.StatusBadRequest)
		return
	}
	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	res.WriteHeader(http.StatusOK)
	res.Write(resp)
}
