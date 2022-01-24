package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/net/context"
)

const (
	_shutdownTimeoutDuration = 5 * time.Second

	_mongoURI  = "mongodb://localhost:27017"
	_serverURI = ":8000"
)

func main() {
	e := echo.New()

	questionMongodb := NewMongo(_mongoURI)
	questionService := NewService(questionMongodb)
	questionHandler := NewHandler(questionService)

	questionHandler.RegisterRoutes(e)

	if err := questionMongodb.Connect(context.Background()); err != nil {
		log.Println("connection could not be established", err)
	}

	go func() {
		if err := e.Start(_serverURI); err != nil {
			log.Println(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), _shutdownTimeoutDuration)
	defer cancel()

	if err := questionMongodb.Disconnect(ctx); err != nil {
		log.Println(err)
	}
	log.Println("Mongo client disconnected")

	if err := e.Close(); err != nil {
		log.Println(err)
	}
	log.Println("Echo server closed")
}
