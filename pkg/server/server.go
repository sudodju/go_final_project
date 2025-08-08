package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/sudodju/go_final_project/pkg/api"
)

func Run() error {
	port := os.Getenv("TODO_PORT")
	http.Handle("/", http.FileServer(http.Dir("web")))
	api.Init()
	return http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
