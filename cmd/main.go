package main

import (
	"context"
	"fmt"
	"github.com/mazxaxz/BroadcastMQ/cmd/broadcast"
	"github.com/mazxaxz/BroadcastMQ/cmd/config"
	"github.com/mazxaxz/BroadcastMQ/cmd/environment"
	"github.com/mazxaxz/BroadcastMQ/pkg/shutdown"
	"github.com/sirupsen/logrus"
	"os"
)

var log = logrus.New()

func main() {
	log.Out = os.Stdout
	fmt.Print(environment.LOGO)

	settings := &environment.Settings{}
	err := settings.LoadSettings()
	if err != nil {
		log.Fatal(err)
	}

	log.SetFormatter(environment.EncodeFormatter(settings.OutputType))
	log.SetLevel(environment.EncodeLogLevel(settings.LogLevel))

	config := &config.Config{}
	if err = config.LoadConfiguration(settings.ConfigPath); err != nil {
		log.WithFields(logrus.Fields{
			"config_file_path": settings.ConfigPath,
		}).Fatal(err)
	}

	if err = config.Validate(); err != nil {
		log.Fatal(err)
	}
	config.FillDefault()

	http := &Http{
		addr:      ":8080",
		liveness:  config.LivenessProbe,
		readiness: config.ReadinessProbe,
		logger:    log,
	}
	go http.ServeHTTP()

	ctx, cancel := context.WithCancel(context.Background())

	go func(c context.Context) {
		bmq := &broadcast.Broadcaster{Config: config.Broadcasts}
		if err = bmq.Initialize(c, log); err != nil {
			log.Fatal(err)
		}
		bmq.Start()
	}(ctx)

	shutdown.GracefulShutdown(cancel, log)
}
