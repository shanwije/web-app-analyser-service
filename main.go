package main

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"os/signal"
	"time"
	"web-app-analyser-service/config"
	"web-app-analyser-service/handlers"
)

var log = logrus.New()

//todo namingConv, docker,
func main() {

	// initiate and make available viper-configurations
	config.SetConfigs()

	readTimeOutDuration := time.Duration(viper.GetInt("server.readTimeout"))
	idleTimeoutDuration := time.Duration(viper.GetInt("server.idleTimeout"))
	timeoutContextDuration := time.Duration(viper.GetInt("server.TimeoutContextDuration"))

	// handlers get initiated in here
	pageAnalytics := handlers.NewPageAnalytics(log)

	serveMux := http.NewServeMux()

	// allocate handler to a pattern
	serveMux.Handle("/page-analytics", pageAnalytics)

	server := &http.Server{
		Addr:        ":" + viper.GetString("server.port"),
		Handler:     serveMux,
		IdleTimeout: idleTimeoutDuration * time.Second,
		ReadTimeout: readTimeOutDuration * time.Second,
	}

	// starting server in a separate goroutine
	go func() {
		log.WithFields(logrus.Fields{
			"port": server.Addr,
		}).Info("Server is starting")
		err := server.ListenAndServe()
		if err != nil {
			log.WithFields(logrus.Fields{
				"error": err,
			}).Fatal("Server closed")
		}
	}()

	// below is to handle graceful shutdown
	signalChannel := make(chan os.Signal)
	signal.Notify(signalChannel, os.Interrupt)
	signal.Notify(signalChannel, os.Kill)

	sig := <-signalChannel
	log.WithFields(logrus.Fields{
		"signal": sig,
	}).Error("Received terminate signal, shutting down")

	// this will keep the server after shutdown signal for a pre-defined duration to complete already received requests
	timeoutContext, _ := context.WithDeadline(context.Background(), time.Now().Add(timeoutContextDuration*time.Second))
	server.Shutdown(timeoutContext)
}
