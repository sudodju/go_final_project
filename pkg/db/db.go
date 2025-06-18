package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

const schema = `CREATE TABLE scheduler (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				date CHAR(8) NOT NULL DEFAULT "",
				title TEXT NOT NULL DEFAULT "",
				comment TEXT NOT NULL DEFAULT "",
				repeat VARCHAR(128) NOT NULL DEFAULT ""
				);
				CREATE INDEX task_date ON scheduler (date);`

func Init(dbFile string) error {
	// Проверяем существует ли файл с именем dbFile("scheduler.db")
	_, err := os.Stat(dbFile)

	// Убеждаемся, что ошибка связана именно с отсутствием файла
	var install bool
	if os.IsNotExist(err) {
		install = true
	}

	// Создаем подключение к БД
	DB, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("Ошибка подключения к БД: %v", err)
	}

	// Если ошибка была связана с отсутствием файла, то создаем новую БД
	if install {

		// Создаем БД
		_, err := DB.Exec(schema)
		if err != nil {
			return fmt.Errorf("Ошибка создания схемы БД: %v", err)
		}
	}
	return nil
}
