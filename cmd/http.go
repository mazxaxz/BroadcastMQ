package main

import (
	"fmt"
	"github.com/mazxaxz/BroadcastMQ/cmd/config"
	"github.com/mazxaxz/BroadcastMQ/pkg/healthchecks"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Http struct {
	addr      string
	liveness  config.Probe
	readiness config.Probe
	logger    *logrus.Logger
}

// ServeHTTP runs http server w(/o) liveness probes
func (h *Http) ServeHTTP() error {
	if h.liveness.Enabled {
		http.HandleFunc(h.liveness.Path, healthchecks.HandleHealth)
	}
	if h.readiness.Enabled {
		http.HandleFunc(h.readiness.Path, healthchecks.HandleReady)
	}

	h.logger.WithFields(logrus.Fields{
		"port": h.addr,
	}).Info("Starting server API Server")

	if err := http.ListenAndServe(h.addr, nil); err != nil {
		return fmt.Errorf("An error occured while starting HTTP server: %v", err)
	}

	return nil
}
