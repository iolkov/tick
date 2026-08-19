package main

import (
	"log"
	"tick/internal/config"
	"tick/internal/info"
	"tick/internal/logger"
	"tick/internal/server"
	"time"
)

var (
	version = ""
	build   = ""
	date    = ""
)

func main() {
	info := &info.AppInfo{
		Version: version,
		Build:   build,
		Date:    date,
		Start:   time.Now(),
	}

	conf, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	log := logger.InitLogger(conf.GetLogLevel())

	srv := server.New(info, conf, log)

	log.Info("Server starting",
		"address", srv.Addr,
	)
	if err := srv.ListenAndServe(); err != nil {
		log.Error("Server failed",
			"error", err,
		)
	}
}
