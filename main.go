package main

import (
	"final_project/pkg/db"
	"final_project/pkg/server"
	"log"
	"os"
)

func main() {
	dbFile := "scheduler.db"
	if envFile := os.Getenv("TODO_DBFILE"); envFile != "" {
		dbFile = envFile
	}
	if err := db.Init(dbFile); err != nil {
		log.Fatal(err)
	}

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}

}
