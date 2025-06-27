package api

import (
	"net/http"
)

func Init() {
	http.HandleFunc("/api/nextdate", nextDayHandler)
	http.HandleFunc("/api/task", taskHandler)
	http.HandleFunc("/api/tasks", tasksHandler)
}

func taskHandler(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		addTaskHandler(res, req)
	case http.MethodGet:
		getTaskHandler(res, req)
	case http.MethodPut:
		updateTaskHandler(res, req)
	}
}
