package db

import (
	"database/sql"
	"fmt"
)

type Task struct {
	ID      string `json:"id,omitempty"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	var id int64

	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err == nil {
		id, err = res.LastInsertId()
	}
	return id, err
}

func Tasks(limit int) ([]*Task, error) {
	rows, err := DB.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?", limit)
	if err != nil {
		return nil, fmt.Errorf("Ошибка запроса SELECT к БД: %v", err)
	}
	defer rows.Close()

	var tasks []*Task

	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, fmt.Errorf("Ошибка при сканировании: %v", err)
		}
		tasks = append(tasks, &task)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if tasks == nil {
		tasks = []*Task{}
	}
	return tasks, nil
}

func GetTask(id string) (*Task, error) {
	query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?"

	row := DB.QueryRow(query, id)

	var task Task

	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("Не найдена запись с таким id: %v", err)
	}
	if err != nil {
		return nil, fmt.Errorf("Не удалось извлечь данные по заданному id: %v", err)
	}
	return &task, nil
}

func UpdateTask(task *Task) error {
	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	res, err := DB.Exec(query, &task.Date, &task.Title, &task.Comment, &task.Repeat, &task.ID)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("Некорректный id для редактирования задачи")
	}
	return nil
}

func UpdateDate(id, next string) error {
	query := `UPDATE scheduler SET date = ? WHERE id = ?`
	res, err := DB.Exec(query, next, id)
	if err != nil {
		return fmt.Errorf("Ошибка при обновлении даты задачи по заданному id: %v", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("Не удалось обновить дату задачи с указанным id: %v", err)
	}

	return nil
}

func DeleteTask(id string) error {
	query := `DELETE FROM scheduler WHERE id = ?`

	res, err := DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("Ошибка при удалении задачи: %v", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("Задача с заданным id не найдена: %v", err)
	}

	return nil
}

func SearchBarGetTasks(search string) ([]*Task, error) {
	// Создаем паттерн поиска FTS5 "ба*" найдет "бассейн"
	searchPattern := search + "*"

	rows, err := DB.Query(`
        SELECT s.id, s.date, s.title, s.comment, s.repeat 
        FROM scheduler s 
        JOIN scheduler_fts fts ON s.id = fts.rowid 
        WHERE scheduler_fts MATCH ? 
        ORDER BY s.date`, searchPattern)

	if err != nil {
		return nil, fmt.Errorf("Ошибка запроса SELECT к FTS: %v", err)
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, fmt.Errorf("Ошибка при сканировании: %v", err)
		}
		tasks = append(tasks, &task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if tasks == nil {
		tasks = []*Task{}
	}
	return tasks, nil
}
