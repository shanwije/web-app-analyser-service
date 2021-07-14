package util

import (
	"net"
	"net/url"
	"strings"
)

func ValidateUrl(str *string) (err string) {
	if *str == "" {
		return "URL param is empty"
	}
	if !isUrl(str) {
		return "Invalid URL"
	}
	return ""
}

func isUrl(str *string) bool {
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

