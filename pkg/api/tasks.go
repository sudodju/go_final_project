package api

import (
	"net/http"

	"github.com/sudodju/go_final_project/pkg/db"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

// Получаем список ближайших задач
func tasksHandler(res http.ResponseWriter, req *http.Request) {
	tasks, err := db.Tasks(10)
	if err != nil {
		writeJsonError(res, err, http.StatusBadRequest)
		return
	}
	writeJson(res, TasksResp{
		Tasks: tasks,
	})
}
