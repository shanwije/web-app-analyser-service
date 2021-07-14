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

//todo comments, logs, tests, packaging, namingConv, docker, error handling
func main() {
	logger := log.New(os.Stdout, "web-app-analyser-service", log.LstdFlags)
	webUrlHandler := handlers.NewPageAnalytics(logger)

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
		log.Println("Server started")
	}()

	signalChannel := make(chan os.Signal)
	signal.Notify(signalChannel, os.Interrupt)
	signal.Notify(signalChannel, os.Kill)

	sig := <-signalChannel
	log.Println("Received terminate signal, shutting down : ", sig)

	timeoutContext, _ := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	server.Shutdown(timeoutContext)
}
