package main

import (
	"fmt"
	"go1f/pkg/server"
	"log"

	"github.com/sudodju/go_final_project/pkg/db"
)

func main() {
	//проверяем существует ли бд
	//создаем новую или подключаемся к уже существующей
	dbFile := "scheduler.db"
	db := db.Init(dbFile)
	if db != nil {
		panic(db)
	}

	fmt.Println("Server is starting")
	err := server.Run()
	if err != nil {
		log.Fatalf("Server not running: %v", err)
	}
}
