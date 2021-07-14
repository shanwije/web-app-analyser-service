package main

import (
	"context"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
	"web-app-analyser-service/config"
	"web-app-analyser-service/handlers"
)

//todo comments, logs, tests, packaging, namingConv, docker, error handling, config file, pointers
func main() {

	config.SetConfigs()

	readTimeOutDuration  := time.Duration(viper.GetInt("server.readTimeout"))
	idleTimeoutDuration  := time.Duration(viper.GetInt("server.idleTimeout"))
	timeoutContextDuration := time.Duration(viper.GetInt("server.TimeoutContextDuration"))

	logger := log.New(os.Stdout, "web-app-analyser-service", log.LstdFlags)
	webUrlHandler := handlers.NewPageAnalytics(logger)

	serveMux := http.NewServeMux()
	serveMux.Handle("/", webUrlHandler)

	server := &http.Server{
		Addr:        ":" + viper.GetString("server.port"),
		Handler:     serveMux,
		IdleTimeout:  idleTimeoutDuration * time.Second,
		ReadTimeout: readTimeOutDuration * time.Second,
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

	timeoutContext, _ := context.WithDeadline(context.Background(), time.Now().Add(timeoutContextDuration*time.Second))
	server.Shutdown(timeoutContext)
}

