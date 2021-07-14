package handlers

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"web-app-analyser-service/collector"
)

type WebUrl struct {
	logger *log.Logger
}

func NewWebUrl(logger *log.Logger) *WebUrl {
	return &WebUrl{logger}
}

func validate(str *string) (err string) {
	if *str == "" {
		return "URL param is empty"
	}
	if !IsUrl(str) {
		return "Invalid URL"
	}
	return ""
}

func IsUrl(str *string) bool {
	url, err := url.ParseRequestURI(*str)
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

	url := request.URL.Query().Get("url")

	if err := validate(&url); err != "" {
		webUrl.logger.Println("Error :", err)
		err := map[string]interface{}{"data": err}
		responseWriter.Header().Set("Content-type", "application/json")
		responseWriter.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(responseWriter).Encode(err)
		return
	}

	data := *collector.GetAppData(url)
	webUrl.logger.Println("data :", data)
	responseWriter.Header().Set("Content-type", "application/json")
	json.NewEncoder(responseWriter).Encode(data)
}
