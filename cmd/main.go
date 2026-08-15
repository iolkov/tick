package main

import (
	"log"
	"tick/internal/server"
)

func main() {
	srv := server.New()

	log.Printf("Server starting on %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
