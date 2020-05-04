package main

import (
	"fmt"
	"github.com/mazxaxz/BroadcastMQ/cmd/config"
	"github.com/mazxaxz/BroadcastMQ/pkg/healthchecks"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Http struct {
}

func (h *Http) ServeHTTP(addr string, liveness *config.Probe, readiness *config.Probe, log *logrus.Logger) error {
	hc := &healthchecks.Healthcheck{}
	if liveness.Enabled {
		http.HandleFunc(liveness.Path, hc.HandleHeath)
	}
	if readiness.Enabled {
		http.HandleFunc(readiness.Path, hc.HandleReady)
	}

	log.WithFields(logrus.Fields{
		"port": addr,
	}).Info("Starting server API Server")

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		return fmt.Errorf("An error occured while starting HTTP server: %v", err)
	}

	return nil
}