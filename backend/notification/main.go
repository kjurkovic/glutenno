package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"notifications/config"
	"notifications/handlers"
	"notifications/mailer"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

var (
	_, b, _, _  = runtime.Caller(0)
	projectRoot = filepath.Dir(b)
)

func main() {
	logger := log.New(os.Stdout, "[notifications-api] ", log.LstdFlags)
	config, _ := config.Load(projectRoot + "/config.yml")

	sender := mailer.SenderFunc(&config.Mailer)
	router := handlers.MessageHandlerFunc(logger, sender)

	sm := mux.NewRouter()
	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/", router.Post)

	// CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"Authorization", "Content-Type", "Access-Control-Allow-Origin", "Origin"},
		AllowedMethods:   []string{"GET", "UPDATE", "PUT", "POST", "DELETE", "OPTIONS"},
		Debug:            true,
	})

	s := &http.Server{
		Addr:         fmt.Sprintf(":%d", config.Server.Port),
		Handler:      c.Handler(sm),
		IdleTimeout:  config.Server.Timeout.Idle * time.Second,
		ReadTimeout:  config.Server.Timeout.Read * time.Second,
		WriteTimeout: config.Server.Timeout.Write * time.Second,
	}

	go func() {
		err := s.ListenAndServe()

		if err != nil {
			logger.Fatal(err)
		}
	}()

	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, os.Interrupt, syscall.SIGTERM)

	sig := <-sigChannel
	logger.Print("Terminate received - shutting down gracefully: ", sig)

	timeoutCtx, _ := context.WithTimeout(context.Background(), config.Server.Timeout.Shutdown*time.Second)
	s.Shutdown(timeoutCtx)

}
