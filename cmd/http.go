package main

import (
	"fmt"
	"github.com/mazxaxz/BroadcastMQ/pkg/healthchecks"
	"net/http"
)

type Http struct {
}

func (h *Http) ServeHTTP(addr string) error {
	hc := &healthchecks.Healthcheck{}
	http.HandleFunc("/_meta/health", hc.HandleHeath)
	http.HandleFunc("/_meta/ready", hc.HandleReady)

	fmt.Printf("Starting server on: http://localhost%s", addr)

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		// TODO: make logging more sophisticated
		fmt.Errorf("An error occured while starting HTTP server: %v", err)
		return err
	}

	return nil
}