package api

import "net/http"

func Init() {
	http.HandleFunc("/api/nextdate", nextDayHandler)
}
