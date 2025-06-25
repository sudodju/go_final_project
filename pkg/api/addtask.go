package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/sudodju/go_final_project/pkg/db"
)

// Проверка правильности формата даты
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
		http.Error(res, "Некорректный метод для добавления задачи", http.StatusBadRequest)
		return
	}

	var task db.Task
	var buf bytes.Buffer

	// читаем тело запроса
	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	// десериализуем json в task
	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	// проверка поля title
	if task.Title == "" {
		http.Error(res, "Поле Title не может быть пустым", http.StatusBadRequest)
		return
	}

	// проверка даты и правильности формата
	if task.Date != "" {
		t, err := checkDate(&task)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
		// проверяем больше ли время t чем время now
		if t.Before(time.Now()) {
			// если задача одноразовая, присваиваем текущее время now
			if task.Repeat == "" {
				task.Date = time.Now().Format(dateFormat)
				// если правило повтора repeat имеется, вычисляем следующую дату для таски и присваиваем время полю Date
			} else {
				next, err := NextDate(time.Now(), task.Date, task.Repeat)
				if err != nil {
					http.Error(res, err.Error(), http.StatusBadRequest)
					return
				}
				task.Date = next
			}
		}
	}
	id, err := db.AddTask(&task)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	writeJson(res, id)
}

/*
		 var band string

	    // берем название группы из параметра `band`
	    band = r.URL.Query().Get("band")

	    // Проверяем POST-запрос или нет
	    if r.Method == http.MethodPost {
	        var artist Artist
	        var buf bytes.Buffer
	        // читаем тело запроса
	        _, err := buf.ReadFrom(r.Body)
	        if err != nil {
	            http.Error(w, err.Error(), http.StatusBadRequest)
	            return
	        }
	        // десериализуем JSON в Artist
	        if err = json.Unmarshal(buf.Bytes(), &artist); err != nil {
	            http.Error(w, err.Error(), http.StatusBadRequest)
	            return
	        }
	        artists[band] = artist
	    } */
