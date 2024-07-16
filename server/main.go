package main

import (
	"Auth-Service/api"
	l "Auth-Service/logger"
	"log"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	conn1, err := grpc.NewClient(":8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("error")
	}

	router := api.NewRouter(conn1)
	err = router.Run(":8080")
	if err != nil {
		logger.Error("error is api get way connection port")
		log.Fatal("error is api get way connection port")
	}

}
