package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/sudodju/go_final_project/pkg/db"
)

// Функция проверки правильности формата даты
func checkDate(task *db.Task) (time.Time, error) {
	var t time.Time
	var err error
	if t, err = time.Parse(dateFormat, task.Date); err != nil {
		return t, fmt.Errorf("Неверно указан формат даты: %v", err)
	}
	return t, nil
}

func addTaskHandler(res http.ResponseWriter, req *http.Request) {

	if req.Method != http.MethodPost {
		writeJsonError(res, fmt.Errorf("Некорректный метод, требуется POST"), http.StatusMethodNotAllowed)
		return
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

	// добавляем таску в бд
	id, err := db.AddTask(&task)
	if err != nil {
		writeJsonError(res, err, http.StatusBadRequest)
		return
	}
	// возвращаем id таски в json
	writeJson(res, map[string]string{"id": fmt.Sprint(id)})
}
