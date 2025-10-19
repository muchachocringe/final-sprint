package server

import (
	"fmt"
	"log"
	"net/http"

	"todo-app/pkg/api"
)

func Start(port int) {
	api.Init()

	addr := fmt.Sprintf(":%d", port)
	log.Printf("Сервер запущен на http://localhost%s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
