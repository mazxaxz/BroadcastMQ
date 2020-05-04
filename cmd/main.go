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

	bmq := &broadcast.Broadcast{}
	bmq.Initialize(config.Broadcasts)

	http := &Http{}
	http.ServeHTTP(":8080", &config.LivenessProbe, &config.ReadinessProbe, log)
}