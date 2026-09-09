package main

import (
	"context"
	"go-chatapp/db"
	"go-chatapp/handler/acc"
	"go-chatapp/handler/contact"
	repository "go-chatapp/repository/generated"
	"go-chatapp/routers"
	"go-chatapp/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	router := gin.Default()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env")
	}

	DB, err := db.Connect(context.Background())
	if err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}
	defer DB.Close()

	queries := repository.New(DB)
	userService := service.NewUserService(queries)
	userHandler := acc.RegisterUserHandler(userService)
	loginHandler := acc.LoginUserHandler(userService)
	routers.SetupUserRoutes(router, userHandler, loginHandler)

	contactService := service.NewContactService(queries)
	listUserHandler := contact.ContactUserListHandler(contactService)
	routers.SetupContactRoutes(router, listUserHandler)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router.Handler(),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server is not runing", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Println("Server Shutdown:", err)
	}
	log.Println("Server exiting")

}
