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

type pageAnalytics struct {
	logger *log.Logger
}

func NewPageAnalytics(logger *log.Logger) *pageAnalytics {
	return &pageAnalytics{logger}
}

func validateUrl(str *string) (err string) {
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

func (webUrl *pageAnalytics) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	webUrl.logger.Println("request came in : ", request.RequestURI)

	url := request.URL.Query().Get("url")

	if err := validateUrl(&url); err != "" {
		webUrl.logger.Println("Error :", err)
		webUrl.serveResponse(responseWriter, err, http.StatusBadRequest)
		return
	}
	webUrl.serveResponse(responseWriter, *collector.GetAppData(url), http.StatusOK)
}

func (webUrl *pageAnalytics) serveResponse(responseWriter http.ResponseWriter, payload interface{}, statusCode int) {
	data := map[string]interface{}{"data": payload}
	responseWriter.Header().Set("Content-type", "application/json")
	responseWriter.WriteHeader(statusCode)
	json.NewEncoder(responseWriter).Encode(data)
}
