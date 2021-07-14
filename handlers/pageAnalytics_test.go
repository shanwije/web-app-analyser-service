package handlers

import "testing"

func TestUrlString_IsUrl(t *testing.T) {
	tests := []struct {
		name      string
		urlString string
		want      bool
	}{
		{"correct URL HTTPS", "https://google.com", true},
		{"correct URL HTTPS", "http://google.com", true},
		{"correct URL with query params", "https://www.youtube.com/watch?v=123", true},
		{"correct IP", "http://172.0.0.1:80", true},
		{"correct IP with query params", "https://172.0.0.1:80/watch?v=123", true},
		{"string", "test123", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsUrl(&tt.urlString); got != tt.want {
				t.Errorf("IsUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}
