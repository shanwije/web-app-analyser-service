package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(responseWriter http.ResponseWriter, request *http.Request) {
		log.Println("request came in : ", request.RequestURI )

		param := request.URL.Query().Get("url")

		if param != "" {
			log.Println(param)
			fmt.Fprintf(responseWriter, "working and the query param is : %s", param)
		} else {
			log.Fatal("url param is not defined in the query")
			http.Error(responseWriter, "url param is not defined in the query", http.StatusBadRequest)
			return
		}
	})

	http.ListenAndServe(":8080", nil)
}
