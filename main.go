package main

import (
	"fmt"
	"go1f/pkg/server"
	"log"

	"github.com/sudodju/go_final_project/pkg/db"
)

func main() {
	// Проверяем существует ли бд
	// Создаем новую или подключаемся к уже существующей
	dbFile := "scheduler.db"
	err := db.Init(dbFile)
	defer db.DB.Close()

	if err != nil {
		panic(err)
	}

	fmt.Println("Запуск сервера")
	errServ := server.Run()
	if errServ != nil {
		log.Fatalf("Сервер не запущен: %v", errServ)
	}
}
