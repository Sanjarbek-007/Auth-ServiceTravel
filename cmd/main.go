package main

import (
	"log"
	"sync"
	router "Auth-Service/api"
	"Auth-Service/api/handlers"
	"Auth-Service/cmd/server"
	"Auth-Service/config"
	l "Auth-Service/logger"
	"Auth-Service/storage/postgres"

	"go.uber.org/zap"
)

var logger *zap.Logger

func initLog() {
	log, err := l.NewLogger()
	if err != nil {
		panic(err)
	}
	logger = log
}

func main() {
	initLog()
	db, err := postgres.ConnectionDb()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	cfg := config.Load()

	router := router.NewRouter(handlers.NewHandler(postgres.NewUserRepository(db), logger))

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		err := router.Run(cfg.HTTPPort)
		if err != nil {
			log.Fatal(err)
		}
	}()

	server.ServerRun(postgres.NewUserRepository(db), &cfg)
	wg.Wait()
}
