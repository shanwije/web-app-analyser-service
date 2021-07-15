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

	if err := util.ValidateUrl(&url); err != "" {
		pageAnalytics.log.WithFields(logrus.Fields{
			"error": err,
		}).Error("Url is invalid")
		serveResponse(responseWriter, err, http.StatusBadRequest)
	}
	serveResponse(responseWriter, *collector.GetAppData(url), http.StatusOK)
}

// sending response back
func serveResponse(responseWriter http.ResponseWriter, payload interface{}, statusCode int) {
	data := map[string]interface{}{"data": payload}
	responseWriter.Header().Set("Content-type", "application/json")
	responseWriter.WriteHeader(statusCode)
	json.NewEncoder(responseWriter).Encode(data)
}
