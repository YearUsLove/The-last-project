package server

import (
	"fmt"
	"net/http"
	"os"

	"final_project/pkg/api"

	
)

var webDir = "./web"

func Start() error {
	port := "7540"
	if envPort := os.Getenv("TODO_PORT"); envPort != "" {
	port = envPort
}

api.InitRoutes()

http.Handle("/", http.FileServer(http.Dir(webDir)))

fmt.Printf("starting server on :%s\n", port)
return http.ListenAndServe(":"+port, nil)
}