package main

import (
	"fmt"
	"github.com/mazxaxz/BroadcastMQ/cmd/broadcast"
	"github.com/mazxaxz/BroadcastMQ/cmd/config"
	"github.com/mazxaxz/BroadcastMQ/cmd/environment"
	"github.com/sirupsen/logrus"
	"os"
)

var log = logrus.New()

func main() {
	log.Out = os.Stdout
	log.SetLevel(logrus.InfoLevel)
	fmt.Print(environment.LOGO)

	settings := &environment.Settings{}
	err := settings.LoadSettings()
	if err != nil {
		log.Fatal(err)
		return
	}

	settings.ConfigureLogging(log)

	config := &config.Config{}
	err = config.LoadConfiguration(settings.ConfigPath)
	if err != nil {
		log.WithFields(logrus.Fields{
			"config_file_path": settings.ConfigPath,
		}).Fatal(err)

		return
	}

	http := &Http{
		addr:      ":8080",
		liveness:  &config.LivenessProbe,
		readiness: &config.ReadinessProbe,
		logger:    log,
	}
	http.ServeHTTP()

	bmq := &broadcast.Broadcast{Config: config.Broadcasts}
	err = bmq.Initialize()
	if err != nil {
		log.Fatal(err)
	}
	bmq.Start()

	// TODO handle shutdown
}