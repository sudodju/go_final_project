package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/sudodju/go_final_project/pkg/api"
)

func Run() error {
	// Если переменная окр. существует, то стартуем на нем, иначе на 7540
	port, ok := os.LookupEnv("TODO_PORT")
	if !ok {
		fmt.Println("Переменной окружения TODO_PORT не существует, сервер запущен на порту 7540")
		port = "7540"
	} else {
		port = os.Getenv("TODO_PORT")
	}

	http.Handle("/", http.FileServer(http.Dir("web")))
	api.Init()
	return http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
