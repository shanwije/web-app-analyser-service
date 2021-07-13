package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
	"web-app-analyser-service/handlers"
)

func main() {
	logger := log.New(os.Stdout, "web-app-analyser-service", log.LstdFlags)
	webUrlHandler := handlers.NewWebUrl(logger)

	serveMux := http.NewServeMux()
	serveMux.Handle("/", webUrlHandler)

	server := &http.Server{
		Addr:        ":8080",
		Handler:     serveMux,
		IdleTimeout: 300 * time.Second,
		ReadTimeout: 120 * time.Second,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			logger.Fatal(err)
		}
	}()

	signalChannel := make(chan os.Signal)
	signal.Notify(signalChannel, os.Interrupt)
	signal.Notify(signalChannel, os.Kill)

	sig := <-signalChannel
	log.Println("Received terminate signal, shutting down : ", sig)

	timeoutContext, _ := context.WithDeadline(context.Background(), time.Now().Add(30*time.Second))
	server.Shutdown(timeoutContext)
}
