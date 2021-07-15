package handlers

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
	"web-app-analyser-service/collector"
	"web-app-analyser-service/util"
)

type pageAnalytics struct {
	log *logrus.Logger
}

func NewPageAnalytics(logger *logrus.Logger) *pageAnalytics {
	return &pageAnalytics{logger}
}

// this implements handler interface
func (pageAnalytics *pageAnalytics) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {

	pageAnalytics.log.WithFields(logrus.Fields{
		"url": request.URL.String(),
	}).Info("Request received")

	url := request.URL.Query().Get("url")
	setupHeader(responseWriter)

	if err := util.ValidateUrl(&url); err != "" {
		pageAnalytics.log.WithFields(logrus.Fields{
			"error": err,
		}).Error("Url is invalid")
		serveResponse(responseWriter, err, http.StatusBadRequest)
	} else {
		serveResponse(responseWriter, *collector.GetAppData(url), http.StatusOK)
	}
}

func setupHeader(responseWriter http.ResponseWriter) {
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	responseWriter.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	responseWriter.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, " +
		"Accept-Encoding, X-CSRF-Token, Authorization")
}

// sending response back
func serveResponse(responseWriter http.ResponseWriter, payload interface{}, statusCode int) {
	data := map[string]interface{}{"data": payload}
	responseWriter.WriteHeader(statusCode)
	json.NewEncoder(responseWriter).Encode(data)
}
