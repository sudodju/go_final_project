package server

import (
	"fmt"
	"net/http"

	"github.com/sudodju/go_final_project/pkg/api"
)

func Run() error {
	port := 7540
	http.Handle("/", http.FileServer(http.Dir("web")))
	api.Init()
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
