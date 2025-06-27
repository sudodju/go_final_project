package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/sudodju/go_final_project/pkg/db"
)

func getTaskHandler(res http.ResponseWriter, req *http.Request) {

	if req.Method != http.MethodGet {
		writeJsonError(res, fmt.Errorf("Некорректный метод, требуется GET"), http.StatusMethodNotAllowed)
	}

	id := req.URL.Query().Get("id")

	if id == "" {
		writeJsonError(res, fmt.Errorf("Не указан идентификатор"), http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJsonError(res, err, http.StatusBadRequest)
	}
	writeJson(res, task)
}

func updateTaskHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPut {
		writeJsonError(res, fmt.Errorf("Некорректный метод, требуется GET"), http.StatusMethodNotAllowed)
	}

	var task db.Task

	// читаем тело запроса в task
	decode := json.NewDecoder(req.Body)
	if err := decode.Decode(&task); err != nil {
		writeJsonError(res, fmt.Errorf("Ошибка декодирования данных"), http.StatusBadRequest)
		return
	}

	// проверка поля title
	if task.Title == "" {
		writeJsonError(res, fmt.Errorf("Поле title не может быть пустым"), http.StatusBadRequest)
		return
	}

	// проверка даты и правильности формата
	now := time.Now()
	if task.Date == "" || task.Date == now.Format(dateFormat) {
		task.Date = now.Format(dateFormat)
	} else {
		t, err := checkDate(&task)
		if err != nil {
			writeJsonError(res, err, http.StatusBadRequest)
			return
		}

		if task.Repeat != "" {
			next, err := NextDate(time.Now(), task.Date, task.Repeat)
			if err != nil {
				writeJsonError(res, err, http.StatusBadRequest)
				return
			}
			task.Date = next
		} else if t.Before(now) {
			// если повторения нет и дата в прошлом, ставим сегодняшнюю
			task.Date = now.Format(dateFormat)
		}
	}

	// Обновляем таску в бд
	err := db.UpdateTask(&task)
	if err != nil {
		writeJsonError(res, err, http.StatusBadRequest)
		return
	}
	writeJson(res, map[string]string{})
}
