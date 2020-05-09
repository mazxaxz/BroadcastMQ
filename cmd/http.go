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
	liveness  *config.Probe
	readiness *config.Probe
	logger    *logrus.Logger
}

func (h *Http) ServeHTTP() error {
	hc := &healthchecks.Healthcheck{}
	if h.liveness.Enabled {
		http.HandleFunc(h.liveness.Path, hc.HandleHeath)
	}
	if h.readiness.Enabled {
		http.HandleFunc(h.readiness.Path, hc.HandleReady)
	}

	h.logger.WithFields(logrus.Fields{
		"port": h.addr,
	}).Info("Starting server API Server")

	err := http.ListenAndServe(h.addr, nil)
	if err != nil {
		return fmt.Errorf("An error occured while starting HTTP server: %v", err)
	}

	return nil
}