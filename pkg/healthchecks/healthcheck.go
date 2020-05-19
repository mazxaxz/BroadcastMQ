package healthchecks

import (
	"encoding/json"
	"net/http"
	"os"
)

// HandleHealth liveness probe handler
func HandleHealth(res http.ResponseWriter, _ *http.Request) {
	self, _ := os.Hostname()
	as := &AppStatus{Host: self}
	as.AddComponent(self, Ok)

	code := http.StatusOK
	if as.IsAnyUnhealthy() {
		code = http.StatusServiceUnavailable
	}

	data, err := json.Marshal(as)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.WriteHeader(code)
	res.Write(data)
}

// HandleReady readiness probe handler
func HandleReady(res http.ResponseWriter, _ *http.Request) {
	self, _ := os.Hostname()
	as := &AppStatus{Host: self}
	as.AddComponent(self, Ok)

	code := http.StatusOK
	if as.IsAnyUnhealthy() {
		code = http.StatusServiceUnavailable
	}

	data, err := json.Marshal(as)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.WriteHeader(code)
	res.Write(data)
}
