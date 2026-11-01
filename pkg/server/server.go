package server

import (
	"log"
	"net/http"
	"os"
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

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(webDir)))

	log.Printf("Server started on :%s", port)
	return http.ListenAndServe(":"+port, mux)
}