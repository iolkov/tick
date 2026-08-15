package main

import (
	"log"
	"tick/internal/info"
	"tick/internal/server"
	"time"
)

var (
	version = "0.0.1-alpha"
	build   = "000000000000"
	date    = time.Now().Format(time.RFC3339)
)

func main() {
	info := &info.AppInfo{
		Version: version,
		Build:   build,
		Date:    date,
		Start:   time.Now(),
	}

	srv := server.New(info)

	log.Printf("Server starting on %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
