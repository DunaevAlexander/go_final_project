package server

import (
	"log"
	"net/http"
	"os"

	"go_final_project/pkg/api"
)

const (
	defaultPort = "7540"
	webDir      = "./web"
)

func Run() error {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = defaultPort
	}

	api.Init()
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	log.Printf("Server started on :%s", port)
	return http.ListenAndServe(":"+port, nil)  // nil → DefaultServeMux
}