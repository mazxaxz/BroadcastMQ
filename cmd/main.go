package main

import (
	"github.com/mazxaxz/BroadcastMQ/pkg/configuration"
	"github.com/mazxaxz/BroadcastMQ/pkg/environment"
)

func main() {
	settings := &environment.Settings{}
	err := settings.LoadSettings()
	if err != nil {
		return
	}

	config := &configuration.Config{}
	err = config.LoadConfiguration(settings.ConfigPath)
	if err != nil {
		return
	}

	// init broadcasting

	http := &Http{}
	http.ServeHTTP(":8080")
}