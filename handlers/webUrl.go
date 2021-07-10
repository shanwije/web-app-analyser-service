package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
)

type WebUrl struct {
	logger *log.Logger
}

type UrlString string

func NewWebUrl(logger *log.Logger) *WebUrl {
	return &WebUrl{logger}
}

func (urlString *UrlString) validate() (err string) {
	if *urlString == "" {
		return "URL param is empty"
	}
	if !urlString.IsUrl() {
		return "Invalid URL"
	}
	return ""
}

func (urlString *UrlString) IsUrl() bool {
	str := string(*urlString)
	url, err := url.ParseRequestURI(str)
	if err != nil {
		return false
	}

	address := net.ParseIP(url.Host)

	if address == nil {
		return strings.Contains(url.Host, ".")
	}
	return true
}

func (webUrl *WebUrl) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	webUrl.logger.Println("request came in : ", request.RequestURI)

	urlString := UrlString(request.URL.Query().Get("url"))

	if err := urlString.validate(); err != "" {
		webUrl.logger.Println("Error :", err)
		err := map[string]interface{}{"data": err}
		responseWriter.Header().Set("Content-type", "application/json")
		responseWriter.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(responseWriter).Encode(err)
		return
	}

	webUrl.logger.Println(urlString)
	fmt.Fprintf(responseWriter, "working and the query param is : %s", urlString)
}
