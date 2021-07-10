package main

import (
	"log"
	"net/http"
	"os"
	"web-app-analyser-service/handlers"
)

func main() {
	logger := log.New(os.Stdout, "web-app-analyser-service", log.LstdFlags)
	webUrlHandler := handlers.NewWebUrl(logger)

	serveMux := http.NewServeMux()
	serveMux.Handle("/", webUrlHandler)

	http.ListenAndServe(":8080", serveMux)

}
