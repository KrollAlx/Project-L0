package app

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
	"project-L0/internal/repository"
	"project-L0/internal/server"
	"project-L0/internal/service"
	"project-L0/internal/transport/http/handler"
	"project-L0/internal/transport/nats"
	"syscall"
)

func Run() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
		return
	}
	DBUser, _ := os.LookupEnv("DB_USER")
	DBPassword, _ := os.LookupEnv("DB_PASSWORD")
	DBName, _ := os.LookupEnv("DB_NAME")
	SSLMode, _ := os.LookupEnv("SSL_MODE")

	db, err := sqlx.Connect("postgres", fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s",
		DBUser, DBPassword, DBName, SSLMode))
	if err != nil {
		log.Println(err)
		return
	}

	repo := repository.New(db)
	s := service.New(repo)
	httpHandler := handler.New(s)
	srv := server.NewHttpServer(httpHandler)

	natsHandler := nats.New(s)
	natsSrv := server.NewNatsServer(natsHandler)

	log.Println("Restore cache")
	err = s.RestoreCache()
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Start nats server")
	go func() {
		err = natsSrv.Run()
		if err != nil {
			log.Println(err)
			return
		}
	}()

	log.Println("Start http server")
	go func() {
		err = srv.Run()
		if err != nil {
			log.Println(err)
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Println("App shutting down")
	err = db.Close()
	if err != nil {
		log.Println(err)
	}
	err = natsSrv.Shutdown()
	if err != nil {
		log.Println(err)
	}
}
