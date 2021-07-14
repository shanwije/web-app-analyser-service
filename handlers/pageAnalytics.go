package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"web-app-analyser-service/collector"
	"web-app-analyser-service/util"
)

type pageAnalytics struct {
	logger *log.Logger
}

func NewPageAnalytics(logger *log.Logger) *pageAnalytics {
	return &pageAnalytics{logger}
}

func (webUrl *pageAnalytics) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	webUrl.logger.Println("request came in : ", request.RequestURI)

	url := request.URL.Query().Get("url")

	if err := util.ValidateUrl(&url); err != "" {
		webUrl.logger.Println("Error :", err)
		serveResponse(responseWriter, err, http.StatusBadRequest)
	}
	serveResponse(responseWriter, *collector.GetAppData(url), http.StatusOK)
}

func serveResponse(responseWriter http.ResponseWriter, payload interface{}, statusCode int) {
	data := map[string]interface{}{"data": payload}
	responseWriter.Header().Set("Content-type", "application/json")
	responseWriter.WriteHeader(statusCode)
	json.NewEncoder(responseWriter).Encode(data)
}
