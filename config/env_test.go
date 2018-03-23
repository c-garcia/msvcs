package config

import (
	"net/url"
	"testing"
)

var getHttpPortTests = []struct {
	url         string
	port        int
	errorRaised bool
}{
	{"https://example.com/", 443, false},
	{"http://example.com/", 80, false},
	{"http://example.com:81/", 81, false},
	{"https://example.com:8000/", 8000, false},
	{"amqp://example.com/", 0, true},
	{"amqp://example.com:3000/", 0, true},
}

func TestGetHttpPort(t *testing.T) {
	for _, tt := range getHttpPortTests {
		u, _ := url.Parse(tt.url)
		port, err := GetHttpPort(u)
		if !tt.errorRaised {
			if port != tt.port {
				t.Errorf("Url: %s expected port %d, got %d", tt.url, tt.port, port)
			}
			if err != nil {
				t.Errorf("Url: %s expected to be valid http. It is not.", tt.url)
			}
		} else {
			if err == nil {
				t.Errorf("Url: %s expected not to be valid http. It is", tt.url)
			}
		}
	}
}
